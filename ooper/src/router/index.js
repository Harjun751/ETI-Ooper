import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import SignUp from "../views/Sign Up.vue";
import Login from "../views/Login.vue";
import NewTrip from "../views/Passenger/New Trip.vue";

const ViewTrips = () => import("../views/Passenger/View Trips.vue")
const UpdateAccount = () => import("../views/Update Account.vue")

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/sign-up",
    name: "sign-up",
    component: SignUp,
  },
  {
    path: "/login",
    name: "login",
    component: Login,
  },
  {
    path: "/new-trip",
    name: "new-trip",
    component: NewTrip,
  },
  {
    path:"/view-trips",
    name:"view-trips",
    component:ViewTrips
  },
  {
    path:"/update-account",
    name:"update-account",
    component:UpdateAccount
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
