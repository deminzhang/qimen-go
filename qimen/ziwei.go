package qimen

// ZiWei 紫微斗数

// ZiWeiGong 紫微宫
var ZiWeiGong = []string{"",
	"命宫", "兄弟", "夫妻", "子女", "财帛", "疾厄", "迁移",
	"交友", //奴仆
	"官禄", "田宅", "福德", "父母",
}

// ZiWeiMingZhu 命主 出生年柱
// 命主是指紫微、廉贞、武曲、贪狼这四颗主星中，与人出生时的年柱相对应的那一颗。
var ZiWeiMingZhu = map[string]string{
	// TODO
}

// ZiWeiShenZhu 身主 出生年支->身主
// 身主是指火星、天相、天梁、天同、文昌、天机中，与人出生时的月柱相对应的那一颗。
var ZiWeiShenZhu = map[string]string{
	"子": "火星", "丑": "天相", "寅": "天梁", "卯": "天同", "辰": "文昌", "巳": "天机",
	"午": "火星", "未": "天相", "申": "天梁", "酉": "天同", "戌": "文昌", "亥": "天机",
}

// ZiWeiRiZhu 日主
// 日主是指太阳、太阴、火星、水星、木星、土星这六颗行星中，与人出生时的日柱相对应的那一颗。
var ZiWeiRiZhu = map[string]string{
	//TODO
}

//庙旺得利平陷闲
//四化 禄权科忌

/*

北斗 主:紫微
正星: 贪狼、巨门、禄存、文曲、廉贞、武曲、破军
助星: 左辅、右弼、擎羊、陀罗
南斗 主:天府
正星: 天梁、天同、天相、七杀、文昌
助星: 天魁、天钺、火星、铃星
中天 主:太阳 太阴
正星 吉: 台辅、封诰、恩光、天贵、天官、天福 三台 八座 龙池 凤阁 天才 天寿 红鸾 天喜 天马 解神 天巫
正星 凶: 天空 地劫 地空 天刑 天姚 天伤 天使 天虚 天哭 孤辰 寡宿 截空 旬空 蜚廉 破碎 天月 阴煞
博士: 博士、力士、青龙、小耗、将军、奏书、飞廉 喜神 病符 大耗 伏兵 官符
长生: 长生、沐浴、冠带、临官、帝旺、衰、病、死、墓、绝、胎、养

十四主星
一.紫微星性属阴土。帝座，尊贵，官禄主。可以解厄，延寿，制化。性格：耳根软，霸气强势，中庸敦厚且高傲。有平、冷漠、尊的特性。
二.贪狼星性属阴水和阳木。主福祸，桃花杀，欲望，演艺。性格：豪爽，贪欲重，反复无常，喜欢文艺和投机。有耻，捣怪，荡的特性。
三.巨门星性属阴水。主是非，化气为暗。口才，地下生意。性格：精明，善辩，苛薄，多口舌是非，言辞锋利。有信，畏怯，邪的特性。
四.廉贞星性属阴火，阳木。化气为囚，代表血光，精密仪器。性格：直率，不修边幅，浮荡，精明算计，刚硬，任性。有寒酸和奸的特性。
五.武曲星性属阴金。财富，财帛主。性格：性刚正直，陷则刚愎自用。刚毅，果断，心地善良，倔强好胜。有义、严酷、雄的特性。
六.破军星性属阴水，阴金。代表杂乱场地，开创变动的工作。性格：胆大，劲头足，积极进取。勇敢和急躁的特性。
七.七杀星性属阴水，阴金。将星，主肃杀。性格：较自大，理智，不轻易服人。勇敢，多疑，多变，侠客风范。有忠贞，凶残，彪悍的特性。
八.天相星性属阴水。官禄主，化气为印。代表吏人和衣食之命。性格：文雅，厚重，奢侈，变色龙。有诚信和冷嘲热讽的特性。
九.天同星性属阳水。有解厄制化作用。福德主。性格：磊落，和蔼，懒散，随和，天真，有和气，多变，温和的特性。
十.天机星性属阴木。兄弟主，益寿，化气为善。性格：机敏，善良，浮夸，计较，善变，应变力强，陷则奸猾。有智慧、多疑、贤德的特性。
十一.天梁星性属阴土。是寿星，化气为荫。福报重，解厄制化作用。性格：老练，敦厚，孤高，默守陈规。有寿，憨直的特性。
十二.天府星性属阳土。代表薪水阶层，土地部门。财帛、田宅主。性格：温和，稳重，怯懦，志高气傲。有礼貌，摆架子，尊贵的特性。
十三.太阳星性属阳火。代表能源，动力，外交。父星，夫星，权贵，官禄主。性格：博爱，刚直，急躁，高调，好动，赤胆忠心。有仁义，激烈，善良的特性。
十四.太阴星性属阴水。代表旅游，运输，美化服务。财帛，田宅主。性格：优雅，柔弱，阴沉，干净，唠叨。有爱意，懒惰，妖的特性。

甲级星：包括紫微、贪狼、巨门、廉贞、武曲、破军等六颗紫微星，以及天府、天相、天梁、天同、七杀、天机、太阳、太阴等八颗天府星。
这些星曜在紫薇斗数中扮演着核心角色，对个人的命运有着重大影响。
乙级星：如天官、天福、天虚、天哭等，这些星曜在命盘中起到辅助甲级主星的作用
丙级星：包括长生十二神、生年博士十二神等，这些星曜在原局命盘中的作用相对较小，但在特定的限运流年盘中可能会变得更加重要。
丁级星：如岁建、龙德等，这些星曜在流年中的作用不容忽视。
戊级星：包括晦气、丧门等，这些星曜通常与不利的影响相关联


*/