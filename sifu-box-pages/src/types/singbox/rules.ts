export const initRuleLabels = [
  { label: "类型", key: "type", allowSort: true, initShow: true },
  { label: "规则", key: "payload", allowSort: true, initShow: true },
  { label: "出站", key: "proxy", allowSort: true, initShow: true },
];
export interface ruleColumns {
  type: string;
  payload: string;
  proxy: string;
  key: string;
}
