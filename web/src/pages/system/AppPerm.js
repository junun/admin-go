import React from 'react';
import {Modal, Checkbox, Row, Col, message, Alert} from 'antd';
import styles from './role.css';
import {connect} from "dva";


@connect(({ loading, user }) => {
  return {
    envAppList: user.envAppList,
    allEnvAppList: user.allEnvAppList,
  }
})

class AppPerm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      temp: [],
      objApp: {},
      roleApp: {},
    }
  }

  componentDidMount() {
    this.props.dispatch({ 
      type: 'user/getRoleEnvApp',
      payload: this.props.rid
    }).then(()=> {
      this.setState({
        roleApp: this.props.envAppList,
      });
    })

    this.props.dispatch({
      type: 'user/getAllEnvApp'
    }).then(()=> {
      var mod = this.props.allEnvAppList
      var obj={}
      for (var x = 0; x < mod.length; x++) {
        var keyname = mod[x].env_id;
        var arr = []
        for (var y = 0; y < mod[x].children.length; y++) {
          arr.push(mod[x].children[y].id)
        }
        obj[keyname] = arr;
      }
      this.setState({
        objApp: obj,
      });
    })
  };

  handleAllCheck = (e, mod) => {
    const checked = e.target.checked;
    var tmp = this.state.objApp;
    var roleAppTmp = this.state.roleApp
    const key = mod.env_id;
    if (tmp.hasOwnProperty(key)) {
      tmp[key].map((item) => {
        if (roleAppTmp.hasOwnProperty(key)) {
          var isIncludes = roleAppTmp[key].includes(item);
          if (checked) {
            if (!isIncludes) {
              // 添加
              roleAppTmp[key].push(item)
            } 
          } else {
            if (isIncludes) {
              // 减少
              var index = roleAppTmp[key].indexOf(item);
              roleAppTmp[key].splice(index, 1);
            }
          }
        } else {
          var tmpArry = []
          roleAppTmp[key] = tmpArry
          roleAppTmp[key].push(item)
        }
      })
    } else {
      message.warning('没有找到对应的全局变量，请刷新后重试！');
      return false;
    }

    this.setState({
      roleApp: roleAppTmp,
    });
  };

  envChecked = (id) => {
    var tmpRoleApp   = this.state.roleApp;
    var tmpObjApp   = this.state.objApp;

    if (tmpRoleApp.hasOwnProperty(id) && tmpObjApp.hasOwnProperty(id)) {
      if (tmpRoleApp[id].length == tmpObjApp[id].length) {
        return true
      } else {
        return false
      }
    } else {
      return false
    }
  }

  isChecked = (env_id, id) => {
    var tmp   = this.state.roleApp;
    if (tmp.hasOwnProperty(env_id)) {
      if (tmp[env_id].includes(id)) {
        return true
      } else {
        return false
      }
    } else {
      return false
    }
  }

  handleAppCheck = (env_id, id) => {
    var tmp   = this.state.roleApp;
    if (tmp.hasOwnProperty(env_id)) {
      if (tmp[env_id].includes(id)) {
        var index = tmp[env_id].indexOf(id);
        tmp[env_id].splice(index, 1)
        this.setState({
          roleApp: tmp,
        });
      } else {
        tmp[env_id].push(id)
        this.setState({
          roleApp: tmp,
        });
      }
    } else {
      var tmpArry = []
      tmp[env_id] = tmpArry
      tmp[env_id].push(id)
      this.setState({
        roleApp: tmp,
      })
    }
  };

  handleAppOk = () => {
    this.setState({loading: true})
    const values = {};
    values.id = this.props.rid;
    values.Apps = this.state.roleApp;
    this.props.dispatch({
      type: 'user/roleAppAdd',
      payload: values,
    }).finally(() => this.setState({loading: false}));
  };

  render() {
    const {allEnvAppList, envAppList} = this.props;
    return (
      <Modal
        visible
        width={1000}
        maskClosable={false}
        title="功能权限设置"
        className={styles.container}
        onCancel={this.props.onCancel}
        confirmLoading={this.state.loading}
        onOk={this.handleAppOk}
      >
        <Alert
          closable
          showIcon
          type="info"
          style={{width: 600, margin: '0 auto 20px', color: '#31708f !important'}}
          message="小提示"
          description={[<div key="1">发布权限仅影响发布功能的发布对象，页面功能权限请在功能权限中设置。</div>,
            <div key="2">如果需要发布权限，请至少设置一个有权限操作的环境，否则无法正常发布。</div>]}
        />
        <table border="1" className={styles.table}>
          <thead>
            <tr>
              <th>环境</th>
              <th>应用</th>
            </tr>
          </thead>
          <tbody>
          { allEnvAppList.map(mod => (
            <tr key={mod.env_id}>
              <td>
                <Checkbox
                  checked={this.envChecked(mod.env_id)}
                  onChange={e => this.handleAllCheck(e, mod)}
                >
                  {mod.env_name}
                </Checkbox>
              </td>
              <td>
                <Row>
                  { mod.children && mod.children.length && mod.children.map(perm => (
                    <Col key={perm.id} span={8}>
                      <Checkbox 
                        checked={this.isChecked(mod.env_id, perm.id)}
                        onChange={() => this.handleAppCheck(mod.env_id, perm.id)}
                      >
                        {perm.app_name} 
                      </Checkbox>
                    </Col>
                  ))}
                </Row>
              </td>
            </tr>
          ))}
          </tbody>
        </table>
      </Modal>
    )
  }
}

export default AppPerm
