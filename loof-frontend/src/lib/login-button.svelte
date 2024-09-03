<script lang="ts">
	import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';
    import { onMount } from 'svelte';

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')), openUsernameModal }: { pocketbase: PocketBase, openUsernameModal: () => void } = $props();
	let isLogin = $state((pocketbase.authStore.model == null) ? false : true);
	let username: string = $state(pocketbase.authStore.model?.username ?? "")

	async function login() {
		try {
			const authData = await pocketbase.collection('users').authWithOAuth2({ provider: 'google' });

			isLogin = true;
			
			pocketbase.collection('users').subscribe(authData.record.id, function (event) {
				switch (event.action) {
					case "update":
						username = event.record.username
						break;
				}
			}, {});
		
			if (authData.record.has_changed_username == false) {
				openUsernameModal();
			}
		} catch (err) {
			isLogin = false;
		}
	}

	async function signout() {
		if (pocketbase.authStore.model != null) {
			pocketbase.collection('users').unsubscribe(pocketbase.authStore.model.id);
		}

		pocketbase.authStore.clear();
		isLogin = false;
	}

	if (pocketbase.authStore.model != null) {
		pocketbase.collection('users').subscribe(pocketbase.authStore.model.id, function (event) {
			console.log("test_user")
			switch (event.action) {
				case "update":
					username = event.record.username
					break;
			}
		}, {});
	}
</script>

<div class="user-actions-container">
	<h2>The Menu:</h2>
	<button>Featured</button>
	<button>Most Bookmarked</button>
	<button disabled={!isLogin}>Your Bookmarks</button>
	{#if isLogin == false}
	<button onclick={login}>Login</button>
	{:else}
	<button onclick={openUsernameModal}>Change Account Settings</button>
	<button onclick={signout}>Sign out of {username}</button>
	{/if}
</div>

<style>
	.user-actions-container {
		display: flex;
		flex-direction: column;

		align-items: flex-end;

		position: absolute;
		z-index: 100;

		background-color: rgb(230, 230, 230);
		border: 5px solid black;

		right: 0;
		bottom: 0;

		margin: 5px;
		padding: 5px;
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