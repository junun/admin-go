import React from 'react';
import {Modal, Checkbox, Row, Col, message, Alert} from 'antd';
import styles from './role.css';

class PagePerm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      temp: [],
      checkedValue: false,
      objPerm: {},
    }
  }

  componentDidMount() {
    const allPerm = this.props.allPerm;
    const rolePerm = this.props.rolePerm;

    this.props.dispatch({ 
      type: 'user/getRolePerms',
      payload: this.props.rid
    }).then(()=> {
      this.setState({
        temp: this.props.rolePerm,
      });
    });

    this.props.dispatch({ 
      type: 'user/getAllPermission' 
    }).then(()=> {
      var mod=this.props.allPerm
      var obj={}
      for (var i = 0; i < mod.length; i++) {
        var key="index"
        var subkey = key + mod[i].id;
        if (mod[i].children && mod[i].children.length) {
          for (var x = 0; x < mod[i].children.length; x++) {
            var keyname = subkey + mod[i].children[x].id;
            var arr = []
            if (mod[i].children[x].children && mod[i].children[x].children.length) {
              for (var y = 0; y < mod[i].children[x].children.length; y++) {
                var tmp = mod[i].children[x].children[y]
                arr.push(tmp.id)
              }
            }
            obj[keyname] = arr;
          }
        }
      }

      this.setState({
        objPerm: obj,
      });
    });
  };

  handleAllCheck = (e, mod, page) => {
    const checked = e.target.checked;
    var tmp = this.state.objPerm;
    var rolePermTemp = this.state.temp;

    const key = "index" +`${mod}` + `${page}`;
    if (tmp.hasOwnProperty(key)) {
      tmp[key].map((item) => {
        var index = rolePermTemp.indexOf(item);
        if (checked) {
          if (index == -1 ) {
            // 添加
            rolePermTemp.push(item)
          }
        } else {
          if (index > -1 ) {
            // 减少
            rolePermTemp.splice(index, 1);
          }
        }
      });
    } else {
      message.warning('没有找到对应的全局变量，请刷新后重试！');
      return false;
    }

    this.setState({
      temp: rolePermTemp,
    });
  };

  handlePermCheck = (id) => {
    var tmp   = this.state.temp;
    var index = tmp.indexOf(id);
    if (index > -1) {
      tmp.splice(index, 1)
      this.setState({
        checkedValue: false,
        temp: tmp,
      });
    } else {
      tmp.push(id)
      this.setState({
        temp: tmp,
        checkedValue: true,
      });
    }
  };

  handlePermissionOk = () => {
    // const keys = this.state.checkedKeys;
    const keys = this.state.temp;
    if (keys.length > 0) {
      const values = {};
      values.id = this.props.rid;
      values.Codes = keys;
      this.props.dispatch({
        type: 'user/rolePermsAdd',
        payload: values,
      });
    } else {
      message.warning('没有内容修改， 请检查。');
      return false;
    }

    this.props.onCancel();
  };

  render() {
    const PermBox = this.PermBox;
    return (
      <Modal
        visible
        width={1000}
        maskClosable={false}
        title="功能权限设置"
        className={styles.container}
        onCancel={this.props.onCancel}
        confirmLoading={this.state.loading}
        onOk={this.handlePermissionOk}
      >
        <Alert
          closable
          showIcon
          type="info"
          style={{width: 600, margin: '0 auto 20px', color: '#31708f !important'}}
          message="小提示"
          description={[<div key="1">功能权限仅影响页面功能，管理发布应用权限请在发布权限中设置。</div>,
            <div key="2">权限更改成功后会强制属于该角色的账户重新登录。</div>]}/>
        <table border="1" className={styles.table}>
          <thead>
          <tr>
            <th>模块</th>
            <th>页面</th>
            <th>功能</th>
          </tr>
          </thead>
          <tbody>
          {this.props.allPerm.map(mod => (
            mod.children && mod.children.length && mod.children.map((page, index) => (
              <tr key={page.id}>
                {index === 0 && <td rowSpan={mod.children.length}>{mod.Name}</td>}
                <td>
                  <Checkbox onChange={e => this.handleAllCheck(e, mod.id, page.id)}>
                    {page.Name}
                  </Checkbox>
                </td>
                <td>
                  <Row>
                    { page.children && page.children.length && page.children.map(perm => (
                      <Col key={perm.id} span={8}>
                        <Checkbox 
                          checked={this.state.temp.includes(perm.id)}
                          onChange={() => this.handlePermCheck(perm.id)}
                        >
                          {perm.Name} 
                        </Checkbox>
                      </Col>
                    ))}
                  </Row>
                </td>
              </tr>
            ))
          ))}
          </tbody>
        </table>
      </Modal>
    )
  }
}

export default PagePerm
