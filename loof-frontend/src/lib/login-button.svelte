<script lang="ts">
	import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')) }: { pocketbase: PocketBase } = $props();

	let authData: RecordAuthResponse<RecordModel> | null = $state(null)
	async function login() {
		authData = await pocketbase.collection('users').authWithOAuth2({ provider: 'google' });
	}

	async function signout() {
		pocketbase.authStore.clear();
		authData = null;
	}
</script>

<div class="node-actions-container">
	{#if authData == null}
	<button onclick={login}>Login</button>
	{:else}
	<button onclick={signout}>Sign out of {authData.record.username}</button>
	{/if}
</div>

<style>
	button {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		height: 100%;
	}
</style>