<template>
  <h1>...oops?</h1>
  <div class="row">
      <input type="text" placeholder="first name" required/>
      <input type="text" placeholder="last name" required/>
  </div>
  <div class="row">
      <input type="text" placeholder="mobile number" required/>
      <input type="email" placeholder="e-mail address" required/>
  </div>
  <div class="row" v-if="isPassenger==false">
      <input type="text" placeholder="identification number" required/>
      <input type="email" placeholder="car license number" required/>
  </div>
  <Button text="update"/>
</template>

<script>
import Button from "../components/button.vue"
import { store } from "../state"
export default {
    components:{Button},
    data(){
        return{
            isPassenger:store.state.isPassenger,
        }
    },
}
</script>

<style scoped>
h1{
    font-size:54px;
    color:var(--dark-yellow);
    font-weight:normal;
    margin-top:150px;
}
.row{
    display:flex;
    justify-content: center;
}
input{
    margin:0 30px 30px 30px;
}
::v-deep(button){
    margin-top:40px;
}
</style>