import React from 'react';
import {Modal, Checkbox, Row, Col, message, Alert} from 'antd';
import styles from './role.css';
import {connect} from "dva";


@connect(({ loading, user }) => {
  return {
    envHostList: user.envHostList,
    allEnvHostList: user.allEnvHostList,
  }
})

class HostPerm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      temp: [],
      objHost: {},
      roleHost: {},
    }
  }

  componentDidMount() {
    this.props.dispatch({ 
      type: 'user/getRoleEnvHost',
      payload: this.props.rid
    }).then(()=> {
      this.setState({
        roleHost: this.props.envHostList,
      });
    })

    this.props.dispatch({
      type: 'user/getAllEnvHost'
    }).then(()=> {
      var mod = this.props.allEnvHostList
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
        objHost: obj,
      });
    })
  };

  handleAllCheck = (e, mod) => {
    const checked = e.target.checked;
    var tmp = this.state.objHost;
    var roleHostTmp = this.state.roleHost
    const key = mod.env_id;
    if (tmp.hasOwnProperty(key)) {
      tmp[key].map((item) => {
        if (roleHostTmp.hasOwnProperty(key)) {
          var isIncludes = roleHostTmp[key].includes(item);
          if (checked) {
            if (!isIncludes) {
              // 添加
              roleHostTmp[key].push(item)
            } 
          } else {
            if (isIncludes) {
              // 减少
              var index = roleHostTmp[key].indexOf(item);
              roleHostTmp[key].splice(index, 1);
            }
          }
        } else {
          var tmpArry = []
          roleHostTmp[key] = tmpArry
          roleHostTmp[key].push(item)
        }
      })
    } else {
      message.warning('没有找到对应的全局变量，请刷新后重试！');
      return false;
    }

    this.setState({
      roleHost: roleHostTmp,
    });
  };

  envChecked = (id) => {
    var tmpRoleHost   = this.state.roleHost;
    var tmpObjHost   = this.state.objHost;

    if (tmpRoleHost.hasOwnProperty(id) && tmpObjHost.hasOwnProperty(id)) {
      if (tmpRoleHost[id].length == tmpObjHost[id].length) {
        return true
      } else {
        return false
      }
    } else {
      return false
    }
  }

  isChecked = (env_id, id) => {
    var tmp   = this.state.roleHost;
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
    var tmp   = this.state.roleHost;
    if (tmp.hasOwnProperty(env_id)) {
      if (tmp[env_id].includes(id)) {
        var index = tmp[env_id].indexOf(id);
        tmp[env_id].splice(index, 1)
        this.setState({
          roleHost: tmp,
        });
      } else {
        tmp[env_id].push(id)
        this.setState({
          roleHost: tmp,
        });
      }
    } else {
      var tmpArry = []
      tmp[env_id] = tmpArry
      tmp[env_id].push(id)
      this.setState({
        roleHost: tmp,
      })
    }
  };

  handleAppOk = () => {
    this.setState({loading: true})
    const values = {};
    values.id = this.props.rid;
    values.Hosts = this.state.roleHost;
    this.props.dispatch({
      type: 'user/roleHostAdd',
      payload: values,
    }).finally(() => this.setState({loading: false}));
  };

  render() {
    const {allEnvHostList, envHostList} = this.props;
    return (
      <Modal
        visible
        width={1000}
        maskClosable={false}
        title="主机权限设置"
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
          description={[<div key="1">主机权限将全局影响属于该角色的用户能够看到的主机。</div>]}
        />
        <table border="1" className={styles.table}>
          <thead>
            <tr>
              <th>环境</th>
              <th>应用</th>
            </tr>
          </thead>
          <tbody>
          { allEnvHostList.map(mod => (
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
                        {perm.host_name} 
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

export default HostPerm
