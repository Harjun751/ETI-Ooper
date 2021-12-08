<template>
  <div v-if="state.isPassenger==null" id="nav">
    <span>ooper</span>
    <router-link class="navigation" to="/">home</router-link>
    <router-link class="navigation" to="/login">login</router-link>
    <router-link class="navigation" to="/sign-up">sign up</router-link>
  </div>
  <div v-else-if="state.isPassenger" id="nav">
    <span>ooper</span>
    <router-link class="navigation" to="/new-trip">new trip</router-link>
    <router-link class="navigation" to="/view-trips">view trips</router-link>
    <router-link class="navigation" to="/update-account">update account</router-link>
    <a class="navigation" @click="signOut">sign out</a>
  </div>
  <div v-else-if="state.isPassenger==false" id="nav">
    <span>ooper</span>
    <router-link class="navigation" to="/trip-management">trip management</router-link>
    <a class="navigation" @click="signOut">sign out</a>
  </div>
  <router-view />
</template>

<script>
import { store } from "./state"
export default {
    data(){
        return{
            state:store.state,
        }
    },
    methods:{
      async signOut(){
        await fetch(process.env.VUE_APP_AUTH_MS_HOST+"/api/v1/authorize",{
            method:"DELETE",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer " +  store.state.jwtAccessToken
            },
            credentials:'include',
            })
            this.$router.push("login")
      }
    },
    async mounted(){
      await fetch(process.env.VUE_APP_AUTH_MS_HOST+"/api/v1/authorize",{
        method:"GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " +  store.state.jwtAccessToken
        },
        credentials:'include',
        })
        .then(async (res)=> await res.json())
        .then((data)=>{
            store.setIsPassenger(data.isPassenger)
            if (data.isPassenger){
              this.$router.push("/new-trip")
            }
            else{
              this.$router.push("/trip-management")
            }
        })
    }
}
</script>


<style>
:root{
 --purple: #2F0B4B;
 --bright-yellow:#D0E322;
 --dark-yellow:#889600
}
#app {
  font-family: "Century Gothic";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  font-size:33px;
}
#nav{
  font-size: 40px;
  display:flex;
  justify-content: center;
  align-items: center;
}
@media screen and (max-width: 1440px) {
  #nav{
    font-size:30px;
  }
  #nav a{
    margin-right:90px !important;
  }
}
#nav span{
  font-size: 54px;
  padding:0;
  margin:0;
  color: var(--bright-yellow);
  margin-right:auto;
}

#nav a {
  color: var(--dark-yellow);
  margin-right:165px;
}
#nav a:last-of-type{
  margin-right:auto;
}

#nav a.router-link-exact-active {
  color: var(--bright-yellow);
}

html{
  background: var(--purple);
}

.navigation {
  overflow: hidden;
  text-decoration: none;
  display:block;
}
.navigation::after {
  content: '';
  position: relative;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 0.1em;
  background-color: var(--dark-yellow);
  opacity: 0;
  transition: opacity 300ms, transform 300ms;
  display: block;
  opacity:1;
  transform: translate3d(-100%, 0, 0);
}
.router-link-active::after{
  background-color: var(--bright-yellow);
}
.navigation:hover::after,
.navigation:focus::after,
.router-link-active::after
{
  transform: translate3d(0, 0, 0);
}

/* input styles */
input{
  all:unset;
  display:block;
  background: none;
  color:var(--bright-yellow);
  text-align: left;
  margin-bottom:30px;
}
input:invalid{
  border-bottom: 3px solid var(--dark-yellow);
}
input:valid{
  border-bottom: 3px solid var(--bright-yellow);
}

/* sweet alert styles */
.custom-swal-modal{
  background:var(--purple) !important;
  font-family: "Century Gothic" !important;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.custom-swal-modal .swal2-title,
.custom-swal-modal .swal2-html-container{
  color:var(--bright-yellow) !important;
}
.custom-swal-icon{
  color:var(--bright-yellow) !important;
  border-color:var(--bright-yellow) !important;
}
.custom-swal-icon .swal2-success-line-tip,
.custom-swal-icon .swal2-success-line-ring,
.custom-swal-icon .swal2-success-line-fix,
.custom-swal-icon .swal2-success-line-long{
  background-color:var(--bright-yellow) !important;
}
.custom-swal-icon .swal2-success-ring{
  border:.25em solid rgba(255,255,255,.2) !important;
}
.custom-swal-content{
  color:var(--bright-yellow) !important;
}
.custom-swal-button{
  padding:0px !important;
  text-align:center;
  background:var(--bright-yellow) !important;
  color:var(--purple) !important;
  font-size:23px !important;
  font-weight:bold !important;
  border-radius: 50px !important;
  width:268px !important;
  height:68px !important;
}
.custom-swal-button:hover{
  background:var(--dark-yellow) !important;
}
</style>
