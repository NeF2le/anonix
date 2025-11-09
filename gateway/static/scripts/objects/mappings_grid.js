class MappingsGrid {
    constructor(toast, { 
        delMappingFn, getMappingsFn
    } = {}) {
        this.toast = toast;
        this.deleteMapping = delMappingFn;
        this.getMappings = getMappingsFn;

        this.defaultColDef = {
            flex: 1,
            filter: true,
            filterParams: { buttons: ["reset"] }
        };

        this.columnDefs = this.buildColumnDefs();
        this.gridOptions = {
            columnDefs: this.columnDefs,
            defaultColDef: this.defaultColDef,
            pagination: true,
            paginationPageSize: 10,
            paginationPageSizeSelector: [10, 20],
        };

        this.api = null;
    }

    async refreshRows() {
        let rows = await this.getMappings();
        this.api.setGridOption("rowData", rows);
    }

    getGridOptions() {
        return this.gridOptions;
    }

    setGridApi(api) {
        this.api = api;
    }

    buildColumnDefs() {
        return [
            {
                field: "id",
                headerName: "Токен",
                filter: "agTextColumnFilter",
                cellClass: "cell--select",
            },
            {
                field: "cipher_text",
                headerName: "Зашифрованный текст",
                filter: "agTextColumnFilter",
            },
            {
                field: "dek_wrapped",
                headerName: "Ключ шифрования",
                filter: "agTextColumnFilter",
            },
            {
                field: "token_ttl",
                headerName: "Время жизни",
                filter: "agTextColumnFilter",
            },
            {
                field: "created_at",
                headerName: "Дата и время создания",
                filter: "agDateColumnFilter",
                cellDataType: "dateTime",
                valueFormatter: (params) => this.formatDateTimeCell(params)
            },
            {
                field: "delete_action",
                headerName: "",
                width: 50,
                pinned: "right",
                sortable: false,
                filter: false,
                cellStyle: {
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center"
                },
                cellRenderer: (params) => this.deleteCellRenderer(params)
            }
        ];
    }

    formatDateTimeCell(params) {
        const v = params.value;
        if (!v) return "";

        const date = new Date(v);
        const dd = String(date.getDate()).padStart(2, "0");
        const mm = String(date.getMonth() + 1).padStart(2, "0");
        const yyyy = date.getFullYear();
        const hh = String(date.getHours()).padStart(2, "0");
        const min = String(date.getMinutes()).padStart(2, "0");
        const ss = String(date.getSeconds()).padStart(2, "0");

        return `${dd}.${mm}.${yyyy} ${hh}:${min}:${ss}`;
    }

    deleteCellRenderer(params) {
        const btn = document.createElement("button");
        btn.className = "btn--delete";
        btn.title = "Удалить";

        btn.innerHTML = `
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20.9997 6.72998C20.9797 6.72998 20.9497 6.72998 20.9197 6.72998C15.6297 6.19998 10.3497 5.99998 5.11967 6.52998L3.07967 6.72998C2.65967 6.76998 2.28967 6.46998 2.24967 6.04998C2.20967 5.62998 2.50967 5.26998 2.91967 5.22998L4.95967 5.02998C10.2797 4.48998 15.6697 4.69998 21.0697 5.22998C21.4797 5.26998 21.7797 5.63998 21.7397 6.04998C21.7097 6.43998 21.3797 6.72998 20.9997 6.72998Z" fill="currentColor"/>
            <path d="M8.49977 5.72C8.45977 5.72 8.41977 5.72 8.36977 5.71C7.96977 5.64 7.68977 5.25 7.75977 4.85L7.97977 3.54C8.13977 2.58 8.35977 1.25 10.6898 1.25H13.3098C15.6498 1.25 15.8698 2.63 16.0198 3.55L16.2398 4.85C16.3098 5.26 16.0298 5.65 15.6298 5.71C15.2198 5.78 14.8298 5.5 14.7698 5.1L14.5498 3.8C14.4098 2.93 14.3798 2.76 13.3198 2.76H10.6998C9.63977 2.76 9.61977 2.9 9.46977 3.79L9.23977 5.09C9.17977 5.46 8.85977 5.72 8.49977 5.72Z" fill="currentColor"/>
            <path d="M15.2104 22.75H8.79039C5.30039 22.75 5.16039 20.82 5.05039 19.26L4.40039 9.19001C4.37039 8.78001 4.69039 8.42001 5.10039 8.39001C5.52039 8.37001 5.87039 8.68001 5.90039 9.09001L6.55039 19.16C6.66039 20.68 6.70039 21.25 8.79039 21.25H15.2104C17.3104 21.25 17.3504 20.68 17.4504 19.16L18.1004 9.09001C18.1304 8.68001 18.4904 8.37001 18.9004 8.39001C19.3104 8.42001 19.6304 8.77001 19.6004 9.19001L18.9504 19.26C18.8404 20.82 18.7004 22.75 15.2104 22.75Z" fill="currentColor"/>
            <path d="M13.6601 17.25H10.3301C9.92008 17.25 9.58008 16.91 9.58008 16.5C9.58008 16.09 9.92008 15.75 10.3301 15.75H13.6601C14.0701 15.75 14.4101 16.09 14.4101 16.5C14.4101 16.91 14.0701 17.25 13.6601 17.25Z" fill="currentColor"/>
            <path d="M14.5 13.25H9.5C9.09 13.25 8.75 12.91 8.75 12.5C8.75 12.09 9.09 11.75 9.5 11.75H14.5C14.91 11.75 15.25 12.09 15.25 12.5C15.25 12.91 14.91 13.25 14.5 13.25Z" fill="currentColor"/>
            </svg>
        `;

        btn.addEventListener("click", (e) => {
            e.stopPropagation();
            this.onDeleteMapping(params);
        });

        return btn;
    }

    async onDeleteMapping(params) {
        const mapping = params.data;
        if (!mapping) return;

        const ok = confirm(`Удалить токен ${mapping.id}?`);
        if (!ok) return;

        try {
            await this.deleteMapping(mapping.id);
            params.api.applyTransaction({ remove: [mapping] });

            this.toast.success("Токен удален");
        } catch (err) {
            console.error(`Failed to delete mapping ${mapping.id}:`, err);
            this.toast.error("Не удалось удалить токен");
        }
    }
}
