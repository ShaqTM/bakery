import Vue from 'vue'
import VueRouter from 'vue-router'


Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'HomeView',
    component: () => import(/* webpackChunkName: "units" */ '../views/HomeView.vue')
  },
  {
    path: '/units',
    name: 'Units',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "units" */ '../views/Units.vue')
  },
  {
    path: '/materials',
    name: 'Materials',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "units" */ '../views/Materials.vue')
  }  
]
const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
