import {Layout, Menu, Icon, Dropdown, Avatar, Button} from 'antd';
import styles from './index.less'

const { Header } = Layout;

const GlobalHeader = ({ title, user, onMenuClick }) => {
  const menu = (
    <Menu selectedKeys={[]} onClick={onMenuClick}>
      <Menu.Item disabled>
        <Icon type="user" />个人中心
      </Menu.Item>
      <Menu.Item disabled>
        <Icon type="setting" />设置
      </Menu.Item>
      {/*<Menu.Item key="triggerError">*/}
        {/*<Icon type="close-circle" />触发报错*/}
      {/*</Menu.Item>*/}
      <Menu.Divider />
      <Menu.Item key="logout">
        <Icon type="logout" />退出登录
      </Menu.Item>
    </Menu>
  );
  return <Header style={{ background: '#fff', padding: '0 20px' }}>
    <span className={styles.pageTitle}>{title}</span>
    <Dropdown overlay={menu}>
      <span className={styles.userMenu}>
        {user.headimgurl && user.headimgurl.length > 0 ?
          <Avatar size="small" src={user.headimgurl}/>
          :
          <Icon type="user"/>
        }
        <span>{user.Nickname}</span>
      </span>
    </Dropdown>
  </Header>;
};

export default GlobalHeader;
