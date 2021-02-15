import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)
export default new Vuex.Store({

  state: {
    count: 0,
    units:[]
//    units:[{id:1,name:"Штука",short_name:"шт."},
//            {id:2,name:"Килограмм",short_name:"кг."},
//        {id:3,name:"Литр",short_name:"л."}
//    ],
//    unit:{id:-1,name:"",short_name:""},
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
      
    },
    updateUnits(state,resp){
      state.units = resp.data
    }
  },
//    updateUnit(state,resp){
//      state.unit.id = resp.data.id
//      state.unit.name = resp.data.name
//      state.unit.short_name = resp.data.short_name
//    }    
//  },
  actions:{
    writeUnit({dispatch},unitData){
      axios({
        url:'/api/writeunit',
        data:unitData,
        method:'POST'
      })
      .then(()=>dispatch('readUnits'))
      .catch(error => console.log(error))
    },
    readUnits({commit}){
      axios({
        url:'/api/readunits',
        method:'GET'
      })
      .then(resp=>commit('updateUnits',resp))
      .catch(err => console.log(err))
    },

    readUnit(id){
      return new Promise((resolve, reject) => {
        axios({
          url:'/api/readunit/',
          method:'GET',
          params:{id:id}
        })
        .then(resp=>resolve(resp))
        .catch(err => {
              console.log(err)
              reject (err)})
      })
    }
  }
})
