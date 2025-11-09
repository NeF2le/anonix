package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	"strings"
)

type HashiCorpAdapter struct {
	client *vaultapi.Client
}

func NewHashiCorpAdapter(client *vaultapi.Client) *HashiCorpAdapter {
	return &HashiCorpAdapter{client: client}
}

func (h *HashiCorpAdapter) GenerateDEK(ctx context.Context, bits int, keyName, context string) ([]byte, []byte, error) {
	resp, err := h.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/datakey/plaintext/%s", keyName), map[string]interface{}{
		"bits":    bits,
		"context": context,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("hashiCorpAdapter.GenerateDEK: failed to generate datakey: %w", err)
	}
	if resp == nil || resp.Data == nil {
		return nil, nil, fmt.Errorf("hashiCorpAdapter.GenerateDEK: empty response")
	}
	wrappedDek := resp.Data["ciphertext"].(string)
	dek, _ := base64.StdEncoding.DecodeString(resp.Data["plaintext"].(string))

	return []byte(wrappedDek), dek, nil
}

func (h *HashiCorpAdapter) UnwrapDEK(ctx context.Context, wrappedDek []byte, keyName string) ([]byte, error) {
	secret, err := h.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/decrypt/%s", keyName), map[string]interface{}{
		"ciphertext": string(wrappedDek),
	})
	if err != nil {
		return nil, fmt.Errorf("hashiCorpAdapter.UnwrapDEK: failed to unwrap dek: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("hashiCorpAdapter.UnwrapDEK: empty response")
	}
	dek, _ := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))

	return dek, nil
}

func (h *HashiCorpAdapter) EncryptDeterministic(ctx context.Context, plaintext []byte, keyName string) ([]byte, error) {
	secret, err := h.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/encrypt/%s", keyName), map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(plaintext),
		"context":   base64.StdEncoding.EncodeToString([]byte("secret")),
	})
	if err != nil {
		return nil, fmt.Errorf("hashiCorpAdapter.EncryptDeterministic: transit encrypt: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("hashiCorpAdapter.EncryptDeterministic: empty response from vault encrypt")
	}

	cipherText, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return nil, fmt.Errorf("hashiCorpAdapter.EncryptDeterministic: ciphertext not found in response")
	}

	return []byte(cipherText), nil
}

func (h *HashiCorpAdapter) DecryptDeterministic(ctx context.Context, ciphertext []byte, keyName string) ([]byte, error) {
	ciphertextStr := string(ciphertext)
	if !strings.HasPrefix(ciphertextStr, "vault:v1:") {
		ciphertextStr = fmt.Sprintf("vault:v1:%s", ciphertextStr)
	}

	path := fmt.Sprintf("transit/decrypt/%s", keyName)
	data := map[string]interface{}{
		"ciphertext": ciphertextStr,
	}

	secret, err := h.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("hashiCorpAdapter.Unwrap: transit decrypt: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("hashiCorpAdapter.Unwrap: empty response from vault decrypt")
	}
	ptB64, ok := secret.Data["plaintext"].(string)
	if !ok {
		return nil, fmt.Errorf("hashiCorpAdapter.Unwrap: plaintext not found in response")
	}
	pt, err := base64.StdEncoding.DecodeString(ptB64)
	if err != nil {
		return nil, fmt.Errorf("hashiCorpAdapter.Unwrap: failed to decode base64 plaintext: %w", err)
	}

	return pt, nil
}
