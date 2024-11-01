<script lang="ts">
	import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';
	let { pocketbase = new PocketBase('http://127.0.0.1:8090'), isLogin = ((pocketbase.authStore.model == null) ? false : true), openUsernameModal, openAccountBookmarksModal, openTopBookmarksModal, openFeaturedBookmarksModal, onClose }: { pocketbase: PocketBase, isLogin: boolean, openUsernameModal: () => void, openAccountBookmarksModal: () => void, openTopBookmarksModal: () => void, openFeaturedBookmarksModal: () => void, onClose: () => void } = $props();
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
			switch (event.action) {
				case "update":
					username = event.record.username
					break;
			}
		}, {});
	}
</script>


<div class="menu-container">
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
	<button onclick={onClose}>Close</button>
</div>

<style>
	.menu-container {
		overflow: hidden;
		overflow-y: auto;

		height: 350px;
		margin-bottom: 30px;
		
		display: flex;
		flex-direction: column;
		align-items: stretch;
		justify-content: space-evenly;
	} 

	h2 {
		font-family: "Jost", "Futura", Arial, Helvetica, sans-serif;
	}
	
	button, a {
		font-family: "Jost", "Futura", Arial, Helvetica, sans-serif;

		background-color: rgb(250, 249, 243);

		border: 2px solid black;
		border-radius: 2px;

		padding: 5px;
	}

	button:hover {
		background-color:  rgb(223, 219, 197);
	}

	button:active {
		color: white;
		background-color: black;
	}
</style>