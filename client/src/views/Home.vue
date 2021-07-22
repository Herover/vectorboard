<template>
  <div class="home">
    <input type="text" placeholder="Write board name here" v-model="newBoardName"/>
    <button v-on:click="createNewBoard()">Create board</button>
    <br/>
    <BoardList v-bind:boards="boards"/>
  </div>
</template>

<script>
// @ is an alias to /src
import BoardList from '@/components/BoardList.vue'
import { API_BASE } from '@/config.js'

export default {
  name: 'Home',
  components: {
    BoardList
  },
  data: () => ({
    boards: [],
    newBoardName: "",
  }),
  mounted: function() {
    this.fetchBoards()
  },
  methods: {
    fetchBoards: function () {
      console.log(API_BASE)
      fetch(API_BASE + '/boards')
        .then(resp => {
          if (resp.ok) {
            return resp.json()
          }
        })
        .then(json => {
          this.boards = json.data;
        });
    },
    createNewBoard: async function () {
      console.log("CREATING");
        const resp = await fetch(`${API_BASE}/boards`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            name: this.newBoardName,
          }),
        });
        const content = await resp.json();

        console.log(content);
        this.fetchBoards();
    },
  }
}
</script>
