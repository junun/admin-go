import React, { Component } from 'react';
import {Icon, Layout, LocaleProvider, Modal} from 'antd';
import UserLayout from './UserLayout';
import Redirect from 'umi/redirect';
import zhCN from 'antd/lib/locale-provider/zh_CN'
import SiderMenu from '@/components/SiderMenu'
import GlobalHeader from '@/components/GlobalHeader'
import GlobalFooter from '@/components/GlobalFooter'
import getPageTitle from '@/utils/getPageTitle'
import {connect} from "dva";

const { Sider, Content } = Layout;

@connect(( { loading, app } ) => {
  return {
    menu: app.menu,
    user: app.user,
    loadingMenu: loading.effects['app/getMenu'],
  };
})
class BasicLayout extends Component {
  menuClick = ({ key }) => {
    const { dispatch } = this.props;
    if (key === 'logout') {
      Modal.confirm({
        title: '确定要退出吗？',
        onOk() {
          dispatch({
            type: 'user/logout'
          });
        },
      });
    }
  };

  menuDidMount = () => {
    this.props.dispatch({
      type: 'app/getMenu',
    });
  };

  render() {
    const { menu, user, loadingMenu, location: { pathname } } = this.props;
    if (pathname === '/user/login') {
      // return <UserLayout>{this.props.children}</UserLayout>
      return <UserLayout></UserLayout>
    }
    if (pathname === '/') {
      if (menu && menu.length > 0) {
        return <Redirect to={menu[0].children[0].Url} />;
      }
    }
    return (
      <LocaleProvider locale={zhCN}>
        <Layout>
            <Sider width={256} style={{ minHeight: '100vh', color: 'white' }} collapsible>
              <div style={{ height: '32px', background: 'rgba(225,225,225,.2)', margin: '16px', textAlign: 'center', padding: '5px',overflow: 'hidden' }}>
                <Icon type="deployment-unit" style={{fontSize: '18px'}} />&nbsp;&nbsp;
                <span style={{fontSize: '16px', }}>Demo 管理平台</span>
              </div>
              <SiderMenu menuData={menu}
                         didMount={this.menuDidMount}
                         loading={loadingMenu}
                         pathname={pathname}/>
            </Sider>
            <Layout>
            <GlobalHeader title={getPageTitle(pathname, menu)}
                          user={user}
                          onMenuClick={this.menuClick} />
            <Content style={{ margin: '24px 16px 0' }}>
              <div style={{ minHeight: 360 }}>
                {this.props.children}
              </div>
            </Content>
            <GlobalFooter copyright="Copyright © 2019 管理平台" />
          </Layout>
        </Layout>
      </LocaleProvider>
    )
  }
}

export default BasicLayout;
