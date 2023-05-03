<template>
    <div class="scale-100">
        <HeaderComponent/>
        <Queue/>
    </div>

</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Queue from "@/components/Queue.vue";
import HeaderComponent from "@/components/HeaderComponent.vue";
import WebSocket from "@/components/Websocket.vue";

@Options({
  components: {
    Queue,HeaderComponent,WebSocket
  },
})
export default class App extends Vue {}
</script>

<style>
</style>
