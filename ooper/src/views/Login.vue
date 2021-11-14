<template>

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
export default {
    components:{Toggle,Button},
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

<style>
#login{
    position:relative;
    float:right;
    top:200px;
    right:50px;
}
input{
    all:unset;
    display:block;
    background: none;
    color:var(--bright-yellow);
    text-align: left;
    margin:0 auto 0 auto;
    margin-bottom:30px;
}
input:invalid{
    border-bottom: 3px solid var(--dark-yellow);
}
input:valid{
    border-bottom: 3px solid var(--bright-yellow);
}
input:last-of-type{
    margin-bottom: 60px;
}
</style>