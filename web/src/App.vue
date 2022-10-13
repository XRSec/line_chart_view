<template>
  <div>
    <div id="main" ref="main"></div>
  </div>
</template>


<script>
export default {
  name: "homePage",
}
</script>

<script setup>
import {onMounted} from "vue";
import * as echarts from "echarts";
import axios from "axios";
import {createToaster} from "@meforma/vue-toaster";

const toaster = createToaster({position: "top-right", duration: 3000, dismissible: false,});

function getData(num) {
  return axios.get('/api/top', {
    params: {
      num: num
    },
    timeout: 99999,
    headers: {
      "Authorization": "Basic dGluZ3poYW5nOnRpbmd6aGFuZw=="
    },
  })
}

let myChart;
// Y轴
const dataOne = [];
// X轴
const time = [];
// 图表选项
let options = {
  title: {text: '网络攻击墙', textStyle: {color: 'black'}},
  dataZoom: [{
    type: "inside",
    realtime: false, //拖动滚动条时是否动态的更新图表数据
    height: 20, //滚动条高度
    start: 80, //滚动条开始位置（共100等份）
    end: 100 //结束位置（共100等份）
  }],
  tooltip: {trigger: 'axis', axisPointer: {type: 'cross', label: {backgroundColor: '#283b56'}}, order: 'seriesAsc',},
  legend: {},
  toolbox: {
    show: true,
    feature: {
      dataView: {readOnly: false},
      magicType: {type: ['bar', 'line', 'stack']},
      restore: {},
      saveAsImage: {}
    }
  },
  xAxis: {type: 'category', data: time, boundaryGap: true}, // 把时间组成的数组接过来，放在x轴上
  yAxis: {type: 'value'},
  series: [{
    data: dataOne,
    type: 'line',
    name: '总流量数',
    markPoint: {data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]},
    markLine: {data: [{type: 'average', name: '平均值'}]}
  },],
};

// 初始化数据
const getAll = async function () {
  await getData(288).then(res => {
    console.log("正在初始化数据")
    time.length = 0
    dataOne.length = 0
    for (let i = 0; i < res.data.data.length; i++) {
      let timeV1 = new Date(Date.parse(res.data.data[i].Time))
      let timeMonth = (timeV1.getMonth() + 1).toString().padStart(2, '0')
      let timeDay = timeV1.getDate().toString().padStart(2, '0')
      let timeHour = timeV1.getHours().toString().padStart(2, '0')
      let timeMinute = timeV1.getMinutes().toString().padStart(2, '0')
      time.unshift(`${timeMonth}-${timeDay} ${timeHour}:${timeMinute}`)
      dataOne.unshift(res.data.data[i].Num)
    }
  }).catch(err => {
    console.log(err)
    toaster.error("初始化数据失败!")
  });
  myChart.hideLoading()
  myChart.setOption(options)
  toaster.info("加载数据完成!")
};
getAll()
// 更新新数据流
setInterval(async function () {
  await getData(1).then(res => {
    let now = new Date(Date.parse(res.data.data[0].Time));  // 获得当前的时间
    let timeMonth = (now.getMonth() + 1).toString().padStart(2, '0')
    let timeDay = now.getDate().toString().padStart(2, '0')
    let timeHour = now.getHours().toString().padStart(2, '0')
    let timeMinute = now.getMinutes().toString().padStart(2, '0')
    if (dataOne[dataOne.length - 1] !== res.data.data[0].Num) {
      time.push(`${timeMonth}-${timeDay} ${timeHour}:${timeMinute}`)
      dataOne.push(res.data.data[0].Num)
      toaster.info("更新一条数据!")
    } else {
      console.log("无数据更新")
      return
    }
    myChart.setOption(options)
  }).catch(err => {
    console.log(err)
    toaster.error("更新数据失败!")
  });
}, 5000)

setInterval(function () {
  getAll()
}, 300000)

onMounted(() => {
  // 初始化图表
  // 新建一个promise对象
  let newPromise = new Promise((resolve) => {
    resolve()
  })
  //然后异步执行echarts的初始化函数
  newPromise.then(() => {
    //	此dom为echarts图标展示dom
    myChart = echarts.init(document.getElementById("main"))
    myChart.showLoading()
    toaster.success("初始化完成!")
  })
})
</script>

<style>
#main {
  margin-top: 20vh;
  width: 100%;
  height: 600px;
}
</style>
