<template>
    <div id="images">
        <Passenger :isPassenger="isPassenger"/>
    </div>
    <div id="login">
        <Toggle @updatePassenger="updatePassenger"/>
        <input type="email" placeholder="e-mail address" required/>
        <input type="password" placeholder="password" required/>
        <Button text="login"/>
    </div>
</template>

<script>
import Toggle from "../components/toggle.vue"
import Button from "../components/button.vue"
import Passenger from "../components/passenger.vue"
export default {
    components:{Toggle,Button,Passenger},
    data(){
        return{
            isPassenger:true
        }
    },
    methods:{
        updatePassenger(passenger){
            this.isPassenger = passenger
        }
    }

}
</script>

<style scoped>
#login{
    position:relative;
    float:right;
    top:200px;
    right:50px;
}
input{
    margin:0 auto 30px auto;
}
input:last-of-type{
    margin-bottom: 60px;
}
#images{
    float:left;
    margin-top:200px;
    margin-left:80px;
}
</style>