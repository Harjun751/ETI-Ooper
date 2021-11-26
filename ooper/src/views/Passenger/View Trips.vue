<template>
    <section>
        <div class="row" v-for="item in data" :key="item.id">
            <span>></span>&nbsp;
            <span>{{item.date}}</span>
            <span class="origin">{{item.PickUp}}</span>
            &nbsp;<span>-</span>&nbsp;
            <span>{{item.DropOff}}</span>
            <span class="time">{{item.time}}</span>
        </div>
    </section>
</template>

<script>
import { store } from "../../state"
export default {
data(){
    return{
        data:[]
    }
},
async mounted(){
    await fetch("http://localhost:5004/api/v1/trips",{
        method:"GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " +  store.state.jwtAccessToken
        },
    })
    .then(async (res)=> await res.json())
    .then((data)=>{
        console.log(data)
        this.data = data
    })
}
}
</script>

<style scoped>
.row{
    border-bottom:3px solid var(--bright-yellow);
    margin:0 100px 0 100px;
    text-align: left;
    color:var(--bright-yellow)
}
section{
    margin-top:150px;
}
.time{
    float:right;
}
.origin{
    margin-left:100px;
}
</style>