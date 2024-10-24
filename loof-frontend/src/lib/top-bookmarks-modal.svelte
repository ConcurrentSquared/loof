<script lang="ts">
	import PocketBase, { type RecordModel } from "pocketbase"
    import { onMount } from "svelte";
    import { getQuotationUsername } from "./node-data.svelte";
    import NodeQuotation from "./node-quotation.svelte";

	let { pocketbase = (new PocketBase('http://127.0.0.1:8090')), onClose }: { pocketbase: PocketBase, onClose: () => void } = $props();
	let bookmarks: RecordModel[] = $state([]);

	onMount(async () => {
		bookmarks = (await pocketbase.collection('most_bookmarked').getList(1, 5, { sort: '-bookmark_count', filter: "bookmark_count > 0", expand: "node, node.author", skipTotal: true })).items;
	})
</script>

<div id="bookmark-container">
	<h2>Most bookmarks:</h2>
	<div class="bookmark-list">
		<ol>
			{#each bookmarks as bookmark}
				<li><NodeQuotation text={bookmark.expand!.node.text} link={"https://www.concurrentsquared.com/loof?x=" + (-bookmark.expand!.node.x_position - 150).toString() + "&y=" + (-bookmark.expand!.node.y_position - 100).toString() + "&zoom=1"} author={bookmark.expand!.node.author} pocketbase={pocketbase}></NodeQuotation></li>
			{/each}
		</ol>
	</div>
	
	<button onclick={onClose}>Close</button>
</div>

<style>
	.bookmark-list {
		overflow: hidden;
		overflow-y: auto;

		height: 350px;
		margin-bottom: 30px;

		background-color: rgb(230, 230, 230);
	} 

	.bookmark-list {
		background-color: inherit;
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