<script lang="ts">
	import PocketBase from "pocketbase";

	let { pocketbase = (new PocketBase('http://127.0.0.1:8090')) }: { pocketbase: PocketBase } = $props();

	let formReference: HTMLElement | undefined = $state();
	let username: string = $state("");

	async function onSubmit(event: SubmitEvent) {
		event.preventDefault();

		console.log(username);
		try {
			const userId: string = pocketbase.authStore.model!.id;
			const record = pocketbase.collection("users").update(userId, { username: username });
		} catch (err) {
			console.log(err);
		}
		(formReference as HTMLElement).parentElement!.remove();
	}
</script>

<form bind:this={formReference} onsubmit={onSubmit}>
	<h1>Account Settings</h1>
	<label for="username">Username:</label>
	<input type="text" name="username" id="username" required bind:value={username}>
	<input type="submit">
</form>