<template>
  <div id="app">
    <nav v-on:click.prevent>
      <router-link
        to="/"
        class="link"
        :class="{'active': $route.fullPath ==='/'}">
        Home
      </router-link>
      <router-link
        to="/signup"
        class="link"
        :class="{'active': $route.fullPath ==='/signup'}">
        Sign up
      </router-link>
      <router-link
        to="/login"
        class="link"
        :class="{'active': $route.fullPath ==='/login'}">
        Login
      </router-link>
    </nav>
    <div class="content">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app',
}
</script>

<style>
#app {
  color: #2c3e50;
  color: #fff;
}

.content {
  margin-top: 40px;
}

nav {
  width: 100vw;
  height: 40px;
  background: rgb(55, 131, 173);
  display: flex;
  flex-flow: row;
}

.link {
  background: rgb(55, 131, 173);
  text-decoration: none;
  color: #fff;
  width: 80px;
  height: 40px;
  display: flex;
  flex-flow: row;
  align-items: center;
  align-content: center;
  justify-content: center;
}

*{
    margin:0;
    padding:0;
}

.active {
  background: #e35885;
}
</style>
