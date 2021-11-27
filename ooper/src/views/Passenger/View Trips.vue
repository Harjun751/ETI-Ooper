<template>
    <section>
        <div class="row" v-for="item in data" :key="item.ID">
            <TripDetails :item="item"/>
        </div>
    </section>
</template>

<script>
import { store } from "../../state"
import TripDetails from "../../components/tripDetails.vue"
export default {
components:{TripDetails},
data(){
    return{
        data:[]
    }
},
methods:{

},
async mounted(){
    var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    await fetch("http://localhost:5004/api/v1/trips",{
        method:"GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer " +  store.state.jwtAccessToken
        },
    })
    .then(async (res)=> await res.json())
    .then((data)=>{
        for (var i in data){
            var date
            if (data[i]["Requested"]["Valid"]){
                date = new Date(data[i]["Requested"]["Time"])
                data[i]["Requested"] = date.getDate() + " " +  months[date.getMonth()] + ", " + date.getFullYear()
                data[i]["Time"] = date.toLocaleTimeString()
            }
            if (data[i]["Start"]["Valid"]){
                date = new Date(data[i]["Start"]["Time"])
                data[i]["Start"] = date.toLocaleDateString() + " " + date.toLocaleTimeString()
            }
            if (data[i]["End"]["Valid"]){
                date = new Date(data[i]["End"]["Time"])
                data[i]["End"] = date.toLocaleDateString() + " " + date.toLocaleTimeString()
            }
        }
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
</style>