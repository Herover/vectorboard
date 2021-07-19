<template>
  <div id="app">
    <BoardList v-bind:boards="boards"/>
  </div>
</template>

<script>
import BoardList from './components/BoardList.vue'

export default {
  name: 'App',
  components: {
    BoardList
  },
  data: () => ({
    boards: [],
  }),
  mounted: function() {
    this.fetchBoards()
  },
  methods: {
    fetchBoards: function () {
      fetch('http://localhost:8080/boards')
        .then(resp => {
          if (resp.ok) {
            return resp.json()
          }
        })
        .then(json => {
          this.boards = json.data;
        });
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
