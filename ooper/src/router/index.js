import { createRouter, createWebHistory } from "vue-router";
import SignUp from "../views/Sign Up.vue";
import Login from "../views/Login.vue";
import NewTrip from "../views/Passenger/New Trip.vue";

const ViewTrips = () => import("../views/Passenger/View Trips.vue")
const UpdateAccount = () => import("../views/Update Account.vue")
const TripManagement = () => import("../views/Driver/Trip Management.vue")

const routes = [
  {
    path: "/sign-up",
    name: "sign-up",
    component: SignUp,
  },
  {
    path: "/login",
    alias: '/',
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
  {
    path:"//trip-management",
    name:"/trip-management",
    component:TripManagement
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
