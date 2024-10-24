<script lang="ts">
	import PocketBase, { type RecordModel } from "pocketbase"
    import { onMount } from "svelte";
    import { getQuotationUsername } from "./node-data.svelte";

	let { text = "", link = "https://www.youtube.com/watch?v=dQw4w9WgXcQ", author = null, pocketbase = (new PocketBase('http://127.0.0.1:8090')) }: { text: string, link: string, author: string | null, pocketbase: PocketBase } = $props();

	let username = $state("")

	onMount(async () => {
		username = await getQuotationUsername(author!, pocketbase)
	})
</script>


<p>"... {text} ..."</p>
<address>- <a href={link}>{username}</a></address>

<style>
	p, address {
		font-family: "EB Garamond", "Futura", Arial, Helvetica, sans-serif;
	}
</style>