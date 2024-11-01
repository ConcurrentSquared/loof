<script lang="ts">
	import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';
	import UpdateAuthorship from './node-actions.svelte'

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')), isLogin = $bindable((pocketbase.authStore.model == null) ? false : true), afterLogin, afterLogout, openUsernameModal, openAccountBookmarksModal, openTopBookmarksModal, openFeaturedBookmarksModal, openMobileMenuModal }: { pocketbase: PocketBase, isLogin: boolean, afterLogin: (userId: string) => Promise<void>, afterLogout: () => Promise<void>, openUsernameModal: () => void, openAccountBookmarksModal: () => void, openTopBookmarksModal: () => void, openFeaturedBookmarksModal: () => void, openMobileMenuModal: () => void } = $props();
	let username: string = $state(pocketbase.authStore.model?.username ?? "")

	async function login() {
		try {
			const authData = await pocketbase.collection('users').authWithOAuth2({ provider: 'google' });
			username = authData.record.username;
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

			await afterLogin(authData.record.id);
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
		await afterLogout();
	}

	if (pocketbase.authStore.model != null) {
		pocketbase.collection('users').subscribe(pocketbase.authStore.model.id, function (event) {
			switch (event.action) {
				case "update":
					username = event.record.username
					break;
			}
		}, {});
	}
</script>

<div class="user-actions-container">
	<h2>Menu:</h2>
	<button onclick={openFeaturedBookmarksModal}>Featured</button>
	<button onclick={openTopBookmarksModal}>Most Bookmarked</button>
	<button disabled={!isLogin} onclick={openAccountBookmarksModal}>Your Bookmarks</button>
	{#if isLogin == false}
	<button onclick={login}>Login</button>
	{:else}
	<button onclick={openUsernameModal}>Change Account Settings</button>
	<button onclick={signout}>Sign out of {username}</button>
	{/if}
	<a href="https://www.concurrentsquared.com/loof/rules">Rules</a>
	<a href="https://www.concurrentsquared.com/loof/privacy">Privacy Policy</a>
	<a href="https://www.concurrentsquared.com/loof/tos">Terms of Service</a>
</div>
<button class="user-actions-button" onclick={openMobileMenuModal}>Menu</button>

<style>
	button, a {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		height: 100%;

		cursor: pointer;

		font-family: inherit;
	}
	
	.user-actions-container {
		display: flex;
		flex-direction: column;

		align-items: flex-end;

		position: absolute;
		z-index: 100;

		background-color: rgb(245, 245, 245);
		border: 5px solid black;

		right: 0;
		bottom: 0;

		margin: 5px;
		padding: 5px;

		font-family: "Jost", "Futura", Arial, Helvetica, sans-serif;
	}

	.user-actions-button {
		position: absolute;
		z-index: 100;

		background-color: rgb(245, 245, 245);
		border: 5px solid black;

		right: 0;
		bottom: 0;

		margin: 5px;
		padding: 5px;

		font-family: "Jost", "Futura", Arial, Helvetica, sans-serif;

		height: 50px;
		width: 50px;
	}

	@media (max-width: 768px) {
        .user-actions-container { display: none; }
    }
	@media (min-width: 768px) {
        .user-actions-button { display: none; }
    }
</style>