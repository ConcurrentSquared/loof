<script lang="ts">
	import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')) }: { pocketbase: PocketBase } = $props();
	let isLogin = $state((pocketbase.authStore.model == null) ? false : true);

	async function login() {
		await pocketbase.collection('users').authWithOAuth2({ provider: 'google' })
			.then((res) => {isLogin = true})
			.catch((err) => {isLogin = false})
	}

	async function signout() {
		pocketbase.authStore.clear();
		isLogin = false;
	}
</script>

<div class="user-actions-container">
	{#if isLogin == false}
	<button onclick={login}>Login</button>
	{:else}
	<button onclick={signout}>Sign out of {pocketbase.authStore.model!.username}</button>
	{/if}
</div>

<style>
	.user-actions-container {
		position: absolute;
		z-index: 100;

		background-color: rgba(230, 230, 230, 0.8);
		right: 0;
		bottom: 0;
	}
	
	button {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		height: 100%;

		cursor: pointer;
	}
</style>