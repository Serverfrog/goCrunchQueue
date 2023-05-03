<template>
    <div>
    </div>
</template>

<script  lang="ts">
import {state} from "@/store/Socket";

export default {
    name: "WebsocketComponent",
    computed: {
        connected(){
            return state.connected;
        }
    }
}

</script>

<style scoped>

</style>