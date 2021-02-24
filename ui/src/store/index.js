import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)
export default new Vuex.Store({

  state: {
    count: 0,
    units:[],
    materials:[],
    recipes:[]
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
    },
    getRecipes:state=>{
      return state.recipes
    },
    getMaterial:state=>id=>{
      return state.materials.find(material => material.id === id);
    },

    getUnit:state=>id=>{
      return state.units.find(unit => unit.id === id);
    }
  },
  mutations: {
    updateUnits(state,resp){
      state.units = resp.data
    },
    updateMaterials(state,resp){
      state.materials = resp.data
    },  
    updateRecipes(state,resp){
      state.recipes = resp.data
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
    
    //Material Price
    writeMaterialPrice({dispatch},priceData){
      axios({
        url:'/api/writematerialprice',
        data:priceData,
        method:'POST'
      })
      .then(()=>dispatch('readMaterials',true))
      .catch(error => console.log(error))
    },

    //Recipes
    writeRecipe({dispatch},recipeData){
      axios({
        url:'/api/writerecipe',
        data:recipeData,
        method:'POST'
      })
      .then(()=>dispatch('readRecipes',true))
      .catch(error => console.log(error))
    },
    readRecipes({commit},price){
      axios({
        url:'/api/readrecipes',
        method:'GET',
        params:{price:price}
      })
      .then(resp=>commit('updateRecipes',resp))
      .catch(err => console.log(err))
    },

    readRecipe({commit},params){
      return new Promise((resolve, reject) => {
        axios({
          url:'/api/readrecipe/',
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
    //Recipe Price
    writeRecipePrice({dispatch},priceData){
      axios({
        url:'/api/writerecipeprice',
        data:priceData,
        method:'POST'
      })
      .then(()=>dispatch('readRecipes',true))
      .catch(error => console.log(error))
    },
  }
})
