import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)
export default new Vuex.Store({

  state: {
    count: 0,
    units:[],
    materials:[]
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
    getMaterials:state=>{
      return state.materials
    }


  },
  mutations: {
    updateUnits(state,resp){
      state.units = resp.data
    },
    updateMaterials(state,resp){
      state.materials = resp.data
    },    
    emptyCommit(){
        
    }
  
  },
//    updateUnit(state,resp){
//      state.unit.id = resp.data.id
//      state.unit.name = resp.data.name
//      state.unit.short_name = resp.data.short_name
//    }    
//  },
  actions:{
    //Units
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

    readUnit({commit},id){
      return new Promise((resolve, reject) => {
        axios({
          url:'/api/readunit/',
          method:'GET',
          params:{id:id}
        })
        .then(resp=>{
            commit('emptyCommit')
            resolve(resp)})
        .catch(err => {
              console.log(err)
              reject (err)})
      })
    },
    //Materials
    writeMaterial({dispatch},materialData){
      axios({
        url:'/api/writematerial',
        data:materialData,
        method:'POST'
      })
      .then(()=>dispatch('readMaterials',true))
      .catch(error => console.log(error))
    },
    readMaterials({commit},price){
      axios({
        url:'/api/readmaterials',
        method:'GET',
        params:{price:price}
      })
      .then(resp=>commit('updateMaterials',resp))
      .catch(err => console.log(err))
    },

    readMaterial({commit},params){
      return new Promise((resolve, reject) => {
        axios({
          url:'/api/readmaterial/',
          method:'GET',
          params:params
        })
        .then(resp=>{
            commit('emptyCommit')
            resolve(resp)})
        .catch(err => {
              console.log(err)
              reject (err)})
      })
    },
    //Price
    writePrice({dispatch},priceData){
      axios({
        url:'/api/writeprice',
        data:priceData,
        method:'POST'
      })
      .then(()=>dispatch('readMaterials',true))
      .catch(error => console.log(error))
    },

  }
})
