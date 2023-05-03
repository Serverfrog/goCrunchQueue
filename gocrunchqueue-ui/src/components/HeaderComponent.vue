<template>
    <div class="sticky top-0 z-30 w-full bg-slate-800">
        <header class="flex flex-row z-50 shadow-xl h-16">
            <div class="basis-16">
                <img src="@/assets/logo.svg" class="h-16">
            </div>
            <div class="basis-64">
                <h1 class="text-2xl w-32 align-baseline">goCrunchQueue</h1>
            </div>
            <div class="basis-16 grow"></div>
            <form class="basis-64 flex flex-row flex-grow row-end-1">
                <div class="basis-64">
                    <input class="rounded-full w-60 bg-slate-600 text-amber-50 h-10 bottom-0 absolute border-amber-50  m-1 pl-4" v-model="itemName"
                           placeholder="Name you want to give it"
                    />
                </div>
                <div class="basis-96">
                    <input class="rounded-full w-80 bg-slate-600 text-amber-50 h-10 bottom-0 absolute border-amber-50 m-1 pl-4" v-model="itemUrl"
                           placeholder="Crunchyroll URL"
                    />
                </div>
                <button class="rounded-full bg-orange-400 basis-32 m-2 row-end-1 float-right" @click="submitQueueItem"
                        v-on:submit.prevent="submitQueueItem" type="button">Add</button>
            </form>


        </header>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "HeaderComponent",
    data() {
        return {
            itemUrl: "",
            itemName: ""
        }
    },
    methods: {
        submitQueueItem() {
            axios.post('/api/add', {
                    "Name": this.itemName,
                    "CrunchyrollUrl": this.itemUrl
                }
            ).then(() => {
                this.itemName = "";
                this.itemUrl = "";
            })
        }
    }
}
</script>

<style scoped>
</style>