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
    },
    getUnit:(state)=>(id)=>state.units.find(unit => unit.id === id)

  },
  mutations: {
    writeUnit(state,payload) {
        if (payload.id!=-1){
            var mUnit = state.units.find(unit => unit.id === payload.id)
            mUnit.name = payload.name
            mUnit.short_name = payload.short_name
        }else{
            state.units.push({id:100,name:payload.name,short_name:payload.short_name})
        }
      
    }
  }
})
