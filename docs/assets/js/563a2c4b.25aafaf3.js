"use strict";(self.webpackChunktfcmt=self.webpackChunktfcmt||[]).push([[609],{5680:(e,t,n)=>{n.d(t,{xA:()=>l,yg:()=>g});var r=n(6540);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=r.createContext({}),p=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},l=function(e){var t=p(e.components);return r.createElement(c.Provider,{value:t},e.children)},u="mdxType",m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},f=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,c=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),u=p(n),f=a,g=u["".concat(c,".").concat(f)]||u[f]||m[f]||o;return n?r.createElement(g,i(i({ref:t},l),{},{components:n})):r.createElement(g,i({ref:t},l))}));function g(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=f;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s[u]="string"==typeof e?e:a,i[1]=s;for(var p=2;p<o;p++)i[p]=n[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}f.displayName="MDXCreateElement"},9407:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>m,frontMatter:()=>o,metadata:()=>s,toc:()=>p});var r=n(8168),a=(n(6540),n(5680));const o={sidebar_position:560},i="Skip posting a comment if there is no change",s={unversionedId:"skip-no-changes",id:"skip-no-changes",title:"Skip posting a comment if there is no change",description:"tfcmt >= v4.4.0 | #773 #774",source:"@site/docs/skip-no-changes.md",sourceDirName:".",slug:"/skip-no-changes",permalink:"/tfcmt/skip-no-changes",draft:!1,editUrl:"https://github.com/suzuki-shunsuke/tfcmt-docs/edit/main/docs/skip-no-changes.md",tags:[],version:"current",sidebarPosition:560,frontMatter:{sidebar_position:560},sidebar:"tutorialSidebar",previous:{title:"Mask sensitive data",permalink:"/tfcmt/mask-sensitive-data"},next:{title:"Output the result to a local file",permalink:"/tfcmt/output-file"}},c={},p=[],l={toc:p},u="wrapper";function m(e){let{components:t,...n}=e;return(0,a.yg)(u,(0,r.A)({},l,n,{components:t,mdxType:"MDXLayout"}),(0,a.yg)("h1",{id:"skip-posting-a-comment-if-there-is-no-change"},"Skip posting a comment if there is no change"),(0,a.yg)("p",null,"tfcmt >= ",(0,a.yg)("a",{parentName:"p",href:"https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v4.4.0"},"v4.4.0")," | ",(0,a.yg)("a",{parentName:"p",href:"https://github.com/suzuki-shunsuke/tfcmt/discussions/773"},"#773")," ",(0,a.yg)("a",{parentName:"p",href:"https://github.com/suzuki-shunsuke/tfcmt/pull/774"},"#774")),(0,a.yg)("p",null,"You can skip posting a comment if there is no change using the command line option ",(0,a.yg)("inlineCode",{parentName:"p"},"-skip-no-changes")," or configuration field ",(0,a.yg)("inlineCode",{parentName:"p"},"disable_comment"),"."),(0,a.yg)("p",null,"e.g."),(0,a.yg)("pre",null,(0,a.yg)("code",{parentName:"pre",className:"language-console"},"$ tfcmt plan -skip-no-changes -- terraform plan\n")),(0,a.yg)("p",null,"tfcmt.yaml"),(0,a.yg)("pre",null,(0,a.yg)("code",{parentName:"pre",className:"language-yaml"},"terraform:\n  plan:\n    when_no_changes:\n      disable_comment: true\n")),(0,a.yg)("p",null,"If the option is set, ",(0,a.yg)("inlineCode",{parentName:"p"},"tfcmt plan")," adds or updates a pull request label but doesn't post a comment if the result of ",(0,a.yg)("inlineCode",{parentName:"p"},"terraform plan")," has no change and no warning."),(0,a.yg)("p",null,"Even if there are no comment, the pull request label lets you know the result.\nThis feature is useful when you want to keep pull request comments clean."))}m.isMDXComponent=!0}}]);