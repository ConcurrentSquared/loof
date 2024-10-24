<script lang="ts">
	import PocketBase, { type RecordModel } from "pocketbase"
    import { onMount } from "svelte";
    import { getQuotationUsername } from "./node-data.svelte";
    import NodeQuotation from "./node-quotation.svelte";

	let { pocketbase = (new PocketBase('http://127.0.0.1:8090')), onClose }: { pocketbase: PocketBase, onClose: () => void } = $props();
	let bookmarks: RecordModel[] = $state([]);

	onMount(async () => {
		bookmarks = await pocketbase.collection('bookmarks').getFullList({ expand: "node, node.author" });
	})

	// TODO: Fix the generation of bookmark links
</script>

<div id="bookmark-container">
	<h2>Your bookmarks:</h2>
	<div class="bookmark-list">
		<ul>
			{#each bookmarks as bookmark}
				<li><NodeQuotation text={bookmark.expand!.node.text} link={"https://www.concurrentsquared.com/loof?x=" + (-bookmark.expand!.node.x_position).toString() + "&y=" + (-bookmark.expand!.node.y_position).toString() + "&zoom=1"} author={bookmark.expand!.node.author} pocketbase={pocketbase}></NodeQuotation></li>
			{/each}
		</ul>
	</div>
	
	<button onclick={onClose}>Close</button>
</div>

<style>
	.bookmark-list {
		overflow: hidden;
		overflow-y: auto;

		height: 350px;
		margin-bottom: 30px;
	} 

	h2 {
		font-family: "Jost", "Futura", Arial, Helvetica, sans-serif;
	}

	button {
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