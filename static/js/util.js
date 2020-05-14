/**
 * 判断是否为空
 * @param {*} data 需要判段的数据
 */
function isEmpty(data){
    return data==""||data==''||data==null||data==[]||data==undefined
}
/**
 * 获取时间
 * @param {*} fmt 时间格式
 * @param {*} time 时间
 */
function getMyTime(fmt, time) {
    fmt = fmt || 'yyyy-MM-dd hh:mm:ss';
    let date = this.isEmpty(time) ? new Date() : time;
    let o = {
        "M+": date.getMonth() + 1,                 //月份
        "d+": date.getDate(),                    //日
        "h+": date.getHours(),                   //小时
        "m+": date.getMinutes(),                 //分
        "s+": date.getSeconds(),                 //秒
        "q+": Math.floor((date.getMonth() + 3) / 3), //季度
        "S": date.getMilliseconds()             //毫秒
    };
    if (/(y+)/.test(fmt))
        fmt = fmt.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (let k in o)
        if (new RegExp("(" + k + ")").test(fmt))
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}
/**
 * 读取cookie 
 * c_name cookie 的key
 */
function getCookieByName(c_name) {
    if (document.cookie.length > 0)     {
       let c_start = document.cookie.indexOf(c_name + "=")           
        if (c_start != -1){ 
            c_start = c_start + c_name.length + 1 
            c_end = document.cookie.indexOf(";", c_start) 
            if (c_end == -1)   
                c_end = document.cookie.length   
                return unescape(document.cookie.substring(c_start, c_end))   
            } 
        } 
    return "" 
}