<template>
  <div id="main" style="width: 100%;height:600px;" ref="main">
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

/*
{"code":200,"data":[{"Num":928,"Time":"2022-10-12 09:24"},{"Num":954,"Time":"2022-10-12 09:36"}]}
*/

function getData(num) {
  return axios.get('/api/top', {
    params: {
      num: num
    },
    timeout: 99999,
    headers: {
      "Authorization": "Basic c2Jjbm06d29haW5pc2hhYmk="
    },
  })
}

// Y轴
const dataOne = (function () {
  return []
})()

// X轴
const time = (function () { // 立即执行函数
  return [];// 存放时间的数组
})();

// 初始化数据
(async function () {
  await getData(119).then(res => {
    console.log("正在初始化数据")
    let now = new Date();  // 获得当前的时间
    for (let i = 0; i < res.data.data.length; i++) {
      let timeV1 = new Date(Date.parse(res.data.data[i].Time))
      let timeMonth = (timeV1.getMonth() + 1).toString().padStart(2, '0')
      let timeDay = timeV1.getDate().toString().padStart(2, '0')
      let timeHour = timeV1.getHours().toString().padStart(2, '0')
      let timeMinute = timeV1.getMinutes().toString().padStart(2, '0')
      time.unshift(`${timeMonth}-${timeDay} ${timeHour}:${timeMinute}`)
      dataOne.unshift(res.data.data[i].Num)
      now = new Date(+now - 720000) // 延迟几秒存储一次？
    }
  }).catch(err => {
    console.log(err)
  })
})()

// 图表选项
let options = {
  title: {
    text: '动态',
    textStyle: {
      color: 'black'
    }
  },
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'cross',
      label: {
        backgroundColor: '#283b56'
      }
    },
    order: 'seriesAsc',
  },
  legend: {},
  //----------------------------我在这！！！
  toolbox: {
    show: true,
    feature: {
      dataView: {readOnly: false},
      magicType: {type: ['bar', 'line', 'stack']},
      restore: {},
      saveAsImage: {}
    }
  },
  // --------------------------------
  xAxis: {
    type: 'category',
    data: time, // 把时间组成的数组接过来，放在x轴上
    boundaryGap: true
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: dataOne,
      type: 'line',
      name: '服务器攻击流量一栏',
      markPoint: {
        data: [
          {type: 'max', name: '最大值'},
          {type: 'min', name: '最小值'}
        ]
      },
      markLine: {
        data: [{type: 'average', name: '平均值'}]
      }
    },
  ],
}


// 初始化图表
function initEcharts() {
  // 新建一个promise对象
  let newPromise = new Promise((resolve) => {
    resolve()
  })
  //然后异步执行echarts的初始化函数
  newPromise.then(() => {
    //	此dom为echarts图标展示dom
    echarts.init(document.getElementById("main")).setOption(options)
  })
}

// 更新新数据流
(async function setInterval() {
  await getData(1).then(res => {
    let now = new Date(Date.parse(res.data.data[0].Time));  // 获得当前的时间
    let timeMonth = (now.getMonth() + 1).toString().padStart(2, '0')
    let timeDay = now.getDate().toString().padStart(2, '0')
    let timeHour = now.getHours().toString().padStart(2, '0')
    let timeMinute = now.getMinutes().toString().padStart(2, '0')
    time.push(`${timeMonth}-${timeDay} ${timeHour}:${timeMinute}`)
    dataOne.push(res.data.data[0].Num)
    initEcharts()
  }).catch(err => {
    console.log(err)
  })
})()

onMounted(() => {
  initEcharts()
})
</script>

<style scoped lang="scss">
</style>
