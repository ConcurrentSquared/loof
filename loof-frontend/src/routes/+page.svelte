<script lang="ts">
	import AccountBookmarksModal from "$lib/account-bookmarks-modal.svelte";
    import FeaturedBookmarksModal from "$lib/featured-bookmarks-modal.svelte";
import LoginButton from "$lib/login-button.svelte";
    import MobileMenuModal from "$lib/mobile-menu-modal.svelte";
	import Modal from "$lib/modal.svelte";
    import TopBookmarksModal from "$lib/top-bookmarks-modal.svelte";
	import TreeViewer from "$lib/tree-viewer.svelte";
    import UsernameModal from "$lib/username-modal.svelte";

	import PocketBase from "pocketbase";

	let pocketbase = $state(new PocketBase('http://127.0.0.1:8090'));
	let isLogin = $state((pocketbase.authStore.model == null) ? false : true)

	enum ModalState {
		None,
		UsernameModal,
		AccountBookmarksModal,
		TopBookmarksModal,
		FeaturedBookmarksModal,
		MobileMenuModal
	}

	let modalState: ModalState = $state(ModalState.None);

	function closeModal() {
		modalState = ModalState.None;
	}

	function openUsernameModal() {
		modalState = ModalState.UsernameModal;
	}

	function openAccountBookmarksModal() {
		modalState = ModalState.AccountBookmarksModal;
	}

	function openTopBookmarksModal() {
		modalState = ModalState.TopBookmarksModal;
	}

	function openFeaturedBookmarksModal() {
		modalState = ModalState.FeaturedBookmarksModal;
	}

	function openMobileMenuModal() {
		modalState = ModalState.MobileMenuModal
	}

	let tree: any | null = $state(null)
	async function afterLogin(userId: string) {
		console.log(userId)
		await tree?.updateAuthorships(userId);
	}

	async function afterLogout() {
		tree?.updateAuthorships("");
	}
</script>

<TreeViewer bind:this={tree} pocketbase={pocketbase} isLogin={isLogin}></TreeViewer>
<LoginButton bind:pocketbase={pocketbase} bind:isLogin={isLogin} afterLogin={afterLogin} afterLogout={afterLogout} openUsernameModal={openUsernameModal} openAccountBookmarksModal={openAccountBookmarksModal} openTopBookmarksModal={openTopBookmarksModal} openFeaturedBookmarksModal={openFeaturedBookmarksModal} openMobileMenuModal={openMobileMenuModal}></LoginButton>

{#if modalState == ModalState.UsernameModal}
	<Modal><UsernameModal pocketbase={pocketbase} afterSubmit={closeModal}></UsernameModal></Modal>
{:else if modalState == ModalState.AccountBookmarksModal}
	<Modal><AccountBookmarksModal pocketbase={pocketbase} onClose={closeModal}></AccountBookmarksModal></Modal>
{:else if modalState == ModalState.TopBookmarksModal}
	<Modal><TopBookmarksModal pocketbase={pocketbase} onClose={closeModal}></TopBookmarksModal></Modal>
{:else if modalState == ModalState.FeaturedBookmarksModal}
	<Modal><FeaturedBookmarksModal pocketbase={pocketbase} onClose={closeModal}></FeaturedBookmarksModal></Modal>
{:else if modalState == ModalState.MobileMenuModal}
	<Modal><MobileMenuModal pocketbase={pocketbase} isLogin={isLogin} openUsernameModal={openUsernameModal} openAccountBookmarksModal={openAccountBookmarksModal} openTopBookmarksModal={openTopBookmarksModal} openFeaturedBookmarksModal={openFeaturedBookmarksModal} onClose={closeModal}></MobileMenuModal></Modal>
{/if}
