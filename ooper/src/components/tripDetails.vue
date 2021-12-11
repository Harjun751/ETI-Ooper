<template>
<div @click="hidden = !hidden">
  <span id="pointer" :class="{ rotate:!hidden }">></span>&nbsp;
    <span>{{item.Requested}}</span>
    <div class="details_container">
        {{item.PickUp}}&nbsp;<span>-</span>&nbsp;{{item.DropOff}}
    </div>
    <span class="time">{{item.Time}}</span>
    <div class="hidden" :class="{ hide: hidden, show: !hidden }">
        <div class="grp">
            <svg width="30" height="30">
            <circle cx="15" cy="15" r="15" />
            </svg>
            <section v-if="item.Start.Valid==false">
                <p>Ride not started</p>
            </section>
            <section v-else>
                <p>Ride Started</p>
                <span>{{item.Start}}</span>
            </section>
        </div>
        <div class="grp">
            <svg version="1.1" id="Capa_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="40px" x="0px" y="0px" viewBox="0 0 498.923 498.923" style="enable-background:new 0 0 498.923 498.923;" xml:space="preserve">
            <g>
                <g>
                    <g>
                        <path d="M249.462,0C151.018,0,70.951,80.106,70.951,178.511c0,92.436,133.617,192.453,172.248,315.948
                            c0.83,2.667,3.322,4.484,6.116,4.465c2.804-0.039,5.256-1.876,6.048-4.563c37.478-126.533,172.6-223.307,172.609-315.869
                            C427.963,80.106,347.886,0,249.462,0z M249.462,313.925c-77.184,0-139.987-62.812-139.987-139.987
                            c0-77.184,62.803-139.987,139.987-139.987c77.165,0,139.977,62.803,139.977,139.987
                            C389.439,251.113,326.626,313.925,249.462,313.925z"/>
                    </g>
                </g>
            </g>
            </svg>
            <section v-if="item.End.Valid==false">
                <p>Ride not ended</p>
            </section>
            <section v-else>
                <p>Ride Ended</p>
                <span>{{item.End}}</span>
            </section>
        </div>
    </div>
</div>
</template>

<script>
export default {
    props:["item"],
    data(){
        return {
            "hidden":true
        }
    }
}
</script>

<style scoped>
div{
    cursor:pointer;
    /* Prevent text selection */
    -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Old versions of Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome, Edge, Opera and Firefox */
}
.time{
    float:right;
}
.hidden{
    transition: height 0.1s;
    overflow: hidden;
}
.hidden.show{
    height:200px;
}
.hidden.hide{
    height:0px;
}
.details_container{
    margin-left:100px;
    max-width:500px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display:inline-block;
}
section{
    width:98vw;
}
.grp{
    display:grid;
    grid-gap:0px;
    grid-template-columns:60px 300px;
    margin-top:10px;
}
.grp p{
    margin:4px;
    margin-bottom:1px;
    grid-column-start:2;
    grid-column-end:2;
    grid-row-start:1;
}
.grp span{
    font-size:22px;
    top:-5px;
    position:relative;
    grid-column-start:2;
    margin-left:5px;
}
.grp svg{
    fill: var(--bright-yellow);
    display: inline;
    grid-area:1 1 2 2;
    position:relative;
    top:20px;
}

#pointer{
    display: inline-block;
    transition: transform 0.1s;
    transform: rotate(0deg);
    -ms-transform: rotate(0deg);	 /*  for IE  */

    /* 	for browsers supporting webkit (such as chrome, firefox, safari etc.). */
    -webkit-transform: rotate(0deg);
}
#pointer.rotate{
    transform: rotate(90deg);
    -ms-transform: rotate(90deg);	 /*  for IE  */

    /* 	for browsers supporting webkit (such as chrome, firefox, safari etc.). */
    -webkit-transform: rotate(90deg);
}
</style>