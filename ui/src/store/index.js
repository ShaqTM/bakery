import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)
export default new Vuex.Store({

//const store = new Vuex.Store({
  state: {
    count: 0,
    units:[{id:1,name:"Штука",short_name:"шт."},
            {id:2,name:"Килограмм",short_name:"кг."},
        {id:3,name:"Литр",short_name:"л."}
    ]
  },
  getters:{
    getUnits:state=>{
        return state.units
    }
  },
  mutations: {
    increment (state) {
      state.count++
    }
  }
})
