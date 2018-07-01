<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <div class="items">
      <div v-for="item in items" class="item" v-on:click="route(item)">
        <h4>{{item.title}}</h4>
        <img v-bind:src="item.src"/>
      </div>
    </div>
  </div>
</template>

<script>
import {server} from '../res';
export default {
  name: 'home',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App',
      items: [],
    }
  },
  methods: {
    route: function(item) {
      this.$router.push({ path: '/p/' + item._id, query: { data: item }})
    }
  },
  beforeMount(){
    // fetch('https://react-eko.herokuapp.com/api/posts')
    //   .then(res => res.json())
    //   .then(res => {
    //     // console.log(res)
    //     this.items = res
    //     console.log(this.items);
    //   })
    //   .catch(e => console.log(e))


    // const token = localStorage.getItem('token');
    // const config = {
    //   headers: {
    //     "Authorization": "Bearer " + token
    //   }
    // }
    //
    // fetch(server + 'restricted', config)
    // .then(res => res.json())
    // .then(res => {
    //   console.log(res);
    //   if(res.message === 'Invalid or expired jwt') {
    //     console.log("Logged out");
    //   }
    // })
    // .catch(e => console.log(e))
 },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.items {
  display: flex;
  flex-flow: row;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
}

.item {
  width: 360px;
  display: flex;
  flex-flow: column;
  align-items: center;
  margin: 6px;
}

.item * {
  margin-top: 10px;
  margin-bottom: 10px;
}

.item img {
  width: 360px;
  height: 360px;
  object-fit: cover;
}

h1, h2 {
  font-weight: normal;
  text-align: center;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
