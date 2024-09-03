<script>
	import LoginButton from "$lib/login-button.svelte";
	import Modal from "$lib/modal.svelte";
	import TreeViewer from "$lib/tree-viewer.svelte";
    import UsernameModal from "$lib/username-modal.svelte";

	import PocketBase from "pocketbase";

	let pocketbase = $state(new PocketBase('http://127.0.0.1:8090'))
</script>

<TreeViewer pocketbase={pocketbase}></TreeViewer>
<LoginButton bind:pocketbase={pocketbase}></LoginButton>
<Modal><UsernameModal pocketbase={pocketbase}></UsernameModal></Modal>