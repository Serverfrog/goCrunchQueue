export class QueueItem {
    private _CrunchyrollUrl: string;
    private _Id: string;
    private _Name: string;

    constructor(CrunchyrollUrl: string, Id: string, Name: string) {
        this._CrunchyrollUrl = CrunchyrollUrl;
        this._Id = Id;
        this._Name = Name;
    }

    set CrunchyrollUrl(value: string) {
        this._CrunchyrollUrl = value;
    }

    get Id(): string {
        return this._Id;
    }

    set Id(value: string) {
        this._Id = value;
    }

    get Name(): string {
        return this._Name;
    }

    set Name(value: string) {
        this._Name = value;
    }

    static getEmpty() {
        return new QueueItem("","","")
    }
}