
const anonymousIconStart = "<svg xmlns='http://www.w3.org/2000/svg' ";
const anonymousIconEnd = " x='0px' y='0px' width='25' height='25'" +
    " viewBox='0 0 1000 1000' enable-background='new 0 0 1000 1000' xml:space='preserve'> " +
    "<g><g transform='translate(0.000000,511.000000) scale(0.100000,-0.100000)'>" +
    "<path d='M3001.6,4371.9c-1065.2-351.8-1955.6-657.5-1982.3-681.8c-60.7-55.8-84.9-140.7-131-436.7c-99.5-657.5-114-" +
    '1790.6-26.7-2438.5c237.8-1810.1,939-3234.3,2135.2-4340.7c264.5-245,693.9-572.6,1006.9-766.7c334.8-206.2,914.7-' +
    '497.4,994.8-497.4c89.8,0,681.8,303.3,1038.5,528.9c446.5,283.9,747.3,526.5,1159.8,939c531.4,528.9,885.6,1014.2,' +
    '1198.6,1637.8C8931.5-616.7,9196,584.4,9196,1943.1c0,662.4-87.3,1531-169.8,1686.3c-14.6,26.7-50.9,63.1-80.1,' +
    '82.5C8873.3,3758,5047,5012.4,4986.3,5010C4959.6,5010,4066.7,4721.3,3001.6,4371.9z M6803.7,3920.6c972.9-320.3,' +
    '1778.5-589.6,1788.2-599.3c12.1-12.1,38.8-155.3,63.1-322.7c60.7-439.2,60.7-1652.3,0-2110.9c-225.6-1681.5-817.7' +
    '-2931-1865.9-3942.8c-458.6-444-909.9-766.7-1514-1084.6l-276.6-145.6l-276.6,145.6c-609,320.3-1065.1,645.4-1516.5,' +
    '1084.6C2157.2-2038.5,1567.6-791.4,1342,875.5c-60.7,448.9-60.7,1691.2,0,2123.1c24.3,167.4,50.9,310.6,63.1,322.' +
    "7c26.7,26.7,3520.6,1176.8,3581.3,1179.2C5013,4500.5,5830.7,4240.8,6803.7,3920.6z'/>" +
    "<path d='M4047.3,2612.8c-143.2-29.1-274.2-131-337.3-262.1c-53.4-106.8-58.2-150.4-65.5-570.2l-9.7-453.7l-106." +
    '8-51c-135.9-67.9-242.6-196.5-288.7-349.4c-53.4-177.1-53.4-2739.3,0-2916.5c46.1-152.9,152.9-281.5,288.7-349.' +
    '4c106.8-50.9,126.2-50.9,1467.9-53.4c1351.5,0,1361.2,0,1470.4,53.4c138.3,67.9,245.1,196.5,291.2,349.4c53.4,177.' +
    '1,53.4,2739.3,0,2916.5c-46.1,152.9-152.9,281.5-288.8,349.4l-106.7,51l-9.7,465.8c-7.3,458.6-7.3,468.3-75.2,579.' +
    '9c-43.7,75.2-106.8,138.3-182,182l-114,67.9l-922,4.9C4552,2629.8,4095.8,2622.5,4047.3,2612.8z M5871.9,1736.' +
    '9v-388.2h-873.5H4125v388.2v388.2h873.5h873.5V1736.9z M6294.1,800.3c19.4-48.5,19.4-2615.6,0-2664.1c-12.1-34-150.' +
    '4-38.8-1295.7-38.8c-1145.2,0-1283.5,4.9-1295.7,38.8c-19.4,48.5-19.4,2615.6,0,2664.1c12.1,34,150.4,38.8,1295.7,' +
    "38.8C6143.7,839.1,6282,834.3,6294.1,800.3z'/>" +
    "<path d='M4663.6,545.5c-366.4-182-509.5-618.7-325.1-992.4l77.6-152.9V-1034c0-662.4,12.1-674.5,582.3-674.5c570.2," +
    '0,582.3,12.1,582.3,674.5v434.3l77.6,152.9c194.1,388.2,29.1,837.1-371.3,1019.1c-75.2,34-160.1,48.5-291.2,48.' +
    '5C4838.3,620.7,4789.8,611,4663.6,545.5z M5163.4,74.8c128.6-106.7,123.7-264.5-12.1-405.2c-63.1-65.5-63.1-70.4-65.' +
    '5-473.1l-2.4-407.6h-84.9h-84.9l-4.9,407.6c-2.4,402.8-2.4,410.1-63.1,473.1c-131,138.3-138.3,276.6-21.8,400.' +
    "3C4908.7,159.7,5059.1,162.2,5163.4,74.8z'/></g></g> </svg>";
const fillGreen = "style='fill: #2389d7'";


var element = new HTMLElement();

export default {
    ANONYMOUS_ICON_ENABLED: anonymousIconStart + fillGreen + anonymousIconEnd,
    ANONYMOUS_ICON_DISABLED: anonymousIconStart + anonymousIconEnd,
};
