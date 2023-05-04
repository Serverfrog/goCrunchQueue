<template>
    <div class="">
        <li class="flex flex-row h-8 bg-emerald-600" v-if="currentItem.Id">
            <div class="grow">
                <p> {{ currentItem.Name}} | {{ currentItemProgress }}</p>
            </div>
            <div class="grow">
                <p> {{ currentItem.CrunchyrollUrl}}</p>
            </div>
            <div class="grow">
                <p> {{ currentItem.Id}}</p>
            </div>
        </li>
        <li v-for="(item) in queueItemList" v-bind:key="item" class="flex flex-row h-8">
            <div class="grow">
                <p> {{ item.Name}}</p>
            </div>
            <div class="grow">
                <p> {{ item.CrunchyrollUrl}}</p>
            </div>
            <div class="grow">
                <p> {{ item.Id}}</p>
            </div>
        </li>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {mapState} from "vuex";

export default defineComponent({
    name: 'QueueComponent',
    props: {
        name: String
    },
    data() {
    },
    mounted() {
        this.$store.commit('updateList')
    },
    methods: {
    },
    computed: {
        ...mapState({
            queueItemList: 'items',
            currentItem: 'currentItem',
            currentItemProgress: 'currentItemProgress'
        })
    }

})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
