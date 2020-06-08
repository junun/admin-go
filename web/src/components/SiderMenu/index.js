import React from 'react';
import {Menu, Icon, Spin, Empty} from 'antd';
import Link from 'umi/link';

const SubMenu = Menu.SubMenu;

class SiderMenu extends React.Component {
  componentDidMount() {
    this.props.didMount();
  }

  render() {
    const { menuData, pathname, loading } = this.props;
    if (loading) {
      return (
        <div style={{width: '100%', textAlign: 'center'}}><Spin /></div>
      );
    }
    if (!menuData || menuData.length === 0) {
      return <Empty />;
    }
    let menuKey;
    const menu = menuData.map(x => {
        if (x.children) {
          const items = x.children.map(item => {
            if (!menuKey) {
              if (pathname !== '/') {
                if (pathname === item.Url) {
                  menuKey = item.id;
                }
              }
            }

            return <Menu.Item key={item.id}>
                <Link to={item.Url}>
                  <Icon type={item.Icon} />{item.Name}
                </Link>
              </Menu.Item>;
          });
          return <SubMenu key={x.id} title={<span><Icon type={x.Icon} /><span>{x.Name}</span></span>}>
            {items}
          </SubMenu>;
        } 
    });
    // console.log(menuKey);
    // let selectedKey;
    // if (menuKey) {
    //   selectedKey = menuKey.toString();
    // } else {
    //   selectedKey = menuData[0].children[0].id.toString();
    // }
    
    // const defaultSelectedKeys = [selectedKey];
    // const defaultOpenKeys = [menuData[0].id.toString()];
    
    return (
      <Menu theme="dark" mode="inline" 
        // defaultSelectedKeys={defaultSelectedKeys}
        // defaultOpenKeys={defaultOpenKeys}
      >
        {menu}
      </Menu>
    );
  }
}

export default SiderMenu;
