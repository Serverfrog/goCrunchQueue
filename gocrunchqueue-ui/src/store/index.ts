import {createStore} from "vuex";
import {QueueItem} from "@/components/entity/QueueItem";
import axios, {AxiosResponse} from 'axios'
import {Event} from "@/store/Socket"

export interface QueueItemList {
    items: Array<QueueItem>
    currentItem: QueueItem
}

export default createStore({
    state: (): QueueItemList => ({
        items: Array<QueueItem>(),
        currentItem: QueueItem.getEmpty()
    }),
    getters: {
        getQueueItemList: state => state.items,
        getCurrentItem: state => state.currentItem,

    },
    mutations: {
        updateList(state) {
            updateList(state)
        },
        eventReceived(state, payload: Array<Event>) {
            let isUpdateable = false;
            for (const event of payload) {
                if (event.Id < 5) {
                    isUpdateable = true;
                }
            }
            if (isUpdateable) {
                updateList(state);
            }

        },
    },
    actions: {},
    modules: {}
})

function updateList(state: /* Vuex store state */ {
    items: QueueItem[];
    currentItem: QueueItem;
}) {
    axios.get('/api/all').then((response: AxiosResponse<Array<QueueItem>, any>) => {
            console.log(response.data)
            state.items.splice(0)
            state.items.push(...response.data)
        }
    )
    axios.get('/api/current').then((response: AxiosResponse<QueueItem, any>) => {
            console.log(response.data);
            state.currentItem = response.data;
        }
    )

}
