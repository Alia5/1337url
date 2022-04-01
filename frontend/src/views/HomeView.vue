<template>
    <AuthComp v-if="!isAuthenticated" />
    <URLShortner v-if="isAuthenticated" />
</template>

<script lang="ts">
import { Vue, Options } from 'vue-class-component';
import AuthComp from '@/components/AuthComp.vue';
import URLShortner from '@/components/URLShortner.vue';
import { ApiHandler } from '@/api';

@Options({
    components: {
        AuthComp,
        URLShortner
    }
})
export default class HomeView extends Vue {
    public isAuthenticated = ApiHandler.isAuthenticated;
}
</script>
