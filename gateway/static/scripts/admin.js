const toast = new Toast();
const apiBase = window.APP_CONFIG.API_BASE_URL;

const tokenizeModalInstance = new Modal(
    document.getElementById('tokenizeOverlay'),
    {
        onSubmit: tokenize,
        toast
    }
);

const detokenizeModalInstance = new Modal(
    document.getElementById('detokenizeOverlay'),
    {
        onSubmit: detokenize,
        toast
    }
);

const ttlSelect = document.getElementById('ttlSelect');
const ttlCustom = document.getElementById('ttlCustom');
if (ttlSelect && ttlCustom) {
    ttlSelect.addEventListener('change', (e) => {
        if (e.target.value === 'custom') {
            ttlCustom.style.display = '';
            ttlCustom.focus();
        } else {
            ttlCustom.style.display = 'none';
            ttlCustom.value = '';
        }
    });
}

async function getSampleRowData() {
    return [
        {
            id: "1789ba2a-a551-4607-80e2-c0a8e1d6abb2",
            cipher_text: "+rULKOuMvakwdJBuNBcxJlTgTfQ6NW+Dwh9GQTVkH+ckbff2rltM8UM3UL6elMel5zSQIGHaB/agCRT+KjjHFWMwBXZ/J4lQOZ8=",
            dek_wrapped: "dmF1bHQ6djE6U1A2ZXZwVzVOTTluU1lna25XVHN3THlra21aWkoyOXlrc3lRNkl6YXc2ME9HRWZaN0YxMVJFWUFJMSt4djRuY3gwVWpTOGtYTG1vQm1PcG4=",
            token_ttl: ttlToHuman({nanos: 3600}),
            created_at: timestampToDate({seconds: 1761832288, nanos: 747647000}),
            deterministic: true,
            reversible: true,
        }
    ]
}

const mappingsGrid = new MappingsGrid(toast, {
    delMappingFn: deleteMapping,
    getMappingsFn: getMappings
});
const mappingsGridOptions = mappingsGrid.getGridOptions();
mappingsGridAPI = agGrid.createGrid(document.querySelector("#mappingList"), mappingsGridOptions);
mappingsGrid.setGridApi(mappingsGridAPI);
mappingsGrid.refreshRows();
// mappingsGrid.setRowData(sampleRowData);

tokenizeBtn = document.getElementById('tokenizeBtn');
detokenizeBtn = document.getElementById('detokenizeBtn');
refreshMappingListBtn = document.getElementById('refreshMappingListBtn');
logoutBtn = document.getElementById('logoutBtn');

tokenizeBtn.addEventListener('click', async () => {
    await tokenizeModalInstance.open()
    await mappingsGrid.refreshRows();
});
detokenizeBtn.addEventListener('click', () => detokenizeModalInstance.open());
refreshMappingListBtn.addEventListener('click', async () => {
    refreshMappingListBtn.classList.add('is-loading');
    await mappingsGrid.refreshRows();
    refreshMappingListBtn.classList.remove('is-loading');
});
logoutBtn.addEventListener('click', () => {
    history.back();
});