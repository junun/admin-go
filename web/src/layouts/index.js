import React, { Component } from 'react';
import {Icon, Layout, LocaleProvider, ConfigProvider, Modal} from 'antd';
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
    notifies: app.notifies,
    loadingMenu: loading.effects['app/getMenu'],
  };
})
class BasicLayout extends Component {
  state = {
    loading: true,
    // notifies: [],
    read: []
  };

  fetch = () => {
    const { dispatch } = this.props;
    this.setState({loading: true});
    dispatch({
      type: 'app/getNotify'
    }).then(() => this.setState({read: []}))
      .finally(() => this.setState({loading: false}))
  }

  handleRead = (e, item) => {
    e.stopPropagation();
    if (this.state.read.indexOf(item.id) === -1) {
      this.state.read.push(item.id);
      this.setState({read: this.state.read});
      const { dispatch } = this.props;
      dispatch({
        type: 'app/patchNotify',
        payload: {ids: item.id.toString()},
      })
    }
  };

  handleReadAll = () => {
    const ids = this.props.notifies.map(x => x.id);
    this.setState({read: ids});
    const { dispatch } = this.props;
    dispatch({
      type: 'app/patchNotify',
      payload: {ids: ids.join(",")},
    })
  };

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

    // this.fetch();
    this.interval = setInterval(this.fetch, 60000)
  };

  render() {
    const {read, loading} = this.state;
    const {notifies, menu, user, loadingMenu, location: { pathname } } = this.props;
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
      <ConfigProvider locale={zhCN}>
        <Layout>
            <Sider width={256} style={{ minHeight: '100vh', color: 'white' }} collapsible>
              <div style={{ height: '32px', background: 'rgba(225,225,225,.2)', margin: '16px', textAlign: 'center', padding: '5px',overflow: 'hidden' }}>
                <Icon type="file-sync" style={{fontSize: '18px'}} />&nbsp;&nbsp;
                <span style={{fontSize: '16px', }}>Spug 管理平台</span>
              </div>
              <SiderMenu menuData={menu}
                         didMount={this.menuDidMount}
                         loading={loadingMenu}
                         pathname={pathname}/>
            </Sider>
            <Layout>
            <GlobalHeader title={getPageTitle(pathname, menu)}
                          user={user}
                          onMenuClick={this.menuClick}
                          loading={loading}
                          notifies={notifies}
                          read={read}
                          handleReadAll={this.handleReadAll}
                          handleRead={this.handleRead}
            />
            <Content style={{ margin: '24px 16px 0' }}>
              <div style={{ minHeight: 360 }}>
                {this.props.children}
              </div>
            </Content>
            <GlobalFooter copyright="Copyright © 2019 管理平台" />
          </Layout>
        </Layout>
      </ConfigProvider>
    )
  }
}

export default BasicLayout;
