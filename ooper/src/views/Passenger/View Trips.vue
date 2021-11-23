<template>
    <section>
        <div class="row" v-for="item in data" :key="item.id">
            <span>></span>&nbsp;
            <span>{{item.date}}</span>
            <span class="origin">{{item.origin}}</span>
            &nbsp;<span>-</span>&nbsp;
            <span>{{item.destination}}</span>
            <span class="time">{{item.time}}</span>
        </div>
    </section>
</template>

<script>
export default {
data(){
    return{
        data:[
            {
                "date":"11 Nov, 2021","origin":"Bukit Batok",
                "destination":"Shuque Secondary","time":"11:38AM",
                "price":"42","id":"asdfasdb12"
            }
            ]
    }
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