import {reactive} from "vue";
import {WebsocketBuilder} from 'websocket-ts';
import {QueueItem} from "@/components/entity/QueueItem";
import store from "@/store/index";

export const state = reactive({
    connected: false,
    events: Array<Event>(),
});



const URL = process.env.NODE_ENV === "production" ? `ws://${window.location.hostname}/ws` : "ws://localhost:80/ws";

export const socket = new WebsocketBuilder(URL)
    .onOpen((i, ev) => {
        console.log("opened");
        state.connected = true;
    })
    .onClose((i, ev) => {
        console.log("closed");
        state.connected = false;
    })
    .onError((i, ev) => {
        console.log("error")
    })
    .onMessage((i, ev: MessageEvent<string>) => {
        const wsEvent = Event.fromJson(JSON.parse(ev.data));
        console.log(`got an ${wsEvent.getReadableEvenId()} event`);
        console.log(`message: ${wsEvent.Message}`);
        state.events.push(wsEvent);
        store.commit('eventReceived',state.events)

    })
    .onRetry((i, ev) => {
        console.log("retry")
    })
    .build();

class EventType {
    public id: number;
    public name: string;

    constructor(id: number, name: string) {
        this.id = id;
        this.name = name;
    }
}

const eventType = new Map<number, EventType>()
eventType.set(1, {id: 1, name: "Added"});
eventType.set(2, {id: 2, name: "Removed"});
eventType.set(3, {id: 3, name: "Process"});
eventType.set(4, {id: 4, name: "Processed"});
eventType.set(5, {id: 5, name: "ErrLogUpdated"});
eventType.set(6, {id: 6, name: "InfoLogUpdated"});

export class Event {
    private _Id: number;
    private _Item: QueueItem;
    private _Message: string;

    constructor(Id: number, Item: QueueItem, Message: string) {
        this._Id = Id;
        this._Item = Item;
        this._Message = Message;
    }

    get Id(): number {
        return this._Id;
    }

    set Id(value: number) {
        this._Id = value;
    }

    get Item(): QueueItem {
        return this._Item;
    }

    set Item(value: QueueItem) {
        this._Item = value;
    }

    get Message(): string {
        return this._Message;
    }

    set Message(value: string) {
        this._Message = value;
    }

    getReadableEvenId(): string {
        return eventType.get(this._Id)!.name;
    }

    static fromJson(json: any): Event {
        const item = QueueItem.getEmpty();
        item.Id = json.Item.Id;
        item.CrunchyrollUrl = json.Item.CrunchyrollUrl;
        item.Name = json.Item.Name;
        return new Event(json.Id, item, json.Message);
    }
}