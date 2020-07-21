import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Row, Popover, Tooltip,
Tag, Col, Button, Popconfirm, Icon, message, Select} from "antd";
import { Link } from 'react-router-dom';
import {connect} from "dva";
import {timeTrans, timeDatetimeTrans, hasPermission} from "@/utils/globalTools"
import moment from 'moment';
import SelectApp from './SelectApp';
import SelectTemplate from './SelectTemplate';
import Approve from './Approve';
import {httpGet, httpPut} from '@/utils/request';
import { stringify } from 'qs';

const FormItem = Form.Item;
const Option = Select.Option;

const checkStatus = {
  '-3': '发布异常',
  '-2': '回滚失败',
  '-1': '已驳回',
  '1': '待审核',
  '2': '审核成功',
  '3': '发布中',
  '4': '回滚待发布',
  '5': '上线成功',
  '6': '回滚成功',
};

@connect(({ loading, deploy, config }) => {
    return {
      deployList: deploy.deployList,
      deployListLoading: loading.effects['deploy/getDeploy'],
      deployLen: deploy.deployLen,
      apiToken: deploy.apiToken,
      projectsList: config.projectsList,
      appTemplateList: config.appTemplateList,
      configEnvList: config.configEnvList,
      gitBranchList: deploy.gitBranchList,
      gitCommitList: deploy.gitCommitList,
    }
 })

class DeployListPage extends React.Component {
  state = {
    addModal: false,
    ext1Visible: false,
    ext2Visible: false,
    approveVisible: false,
    selectAppVisible: false,
    selectTemplateVisible: false,
    editCacheData: {},
    appId: 0,
    gitShow: 0,
    disableValues: false,
    showTemplate: false,
    undoTemplate: false,
    id: 0,
    branchLoading: false,
    commitLoading: false,
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'deploy/getDeploy',
      payload: {
        page: 1,
        pageSize: 10,
      }
    });

    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });

    dispatch({
      type: 'config/getAppTemplate',
      payload: {
        aid: 0,
      }
    });

    dispatch({
      type: 'config/getProject',
      payload: {
        page: 1,
        pageSize: 999,
        active: 1,
      }
    });
  }
  
  showAddModal = () => {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    })
    dispatch({
      type: 'config/getProject',
      payload: {
        page: 1,
        pageSize: 999,
        active: 1,
      }
    })
    this.setState({ 
      selectAppVisible: true,
      addModal: true,
    })
  }

  selectCancel = () => {
    this.setState({
      selectAppVisible: false,
      selectTemplateVisible:false,
      addModal: false,
      editCacheData: {},
    })
  }

  handleApprove = (id) => {
    this.setState({
      approveVisible: true,
      id: id,
    })
  }

  approveCanael = () => {
    this.setState({
      approveVisible: false,
    })
  }

  handleEx1Cancel = () => {
    this.setState({
      visible: false,
      gitShow: 0,
      addModal: false,
      disableValues: false,
      editCacheData: {},
    });
    // 清理 branchList
    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/cleanBranchList',
    });
  };

  switchTemplate = () => {
    this.setState({ 
      showTemplate: !this.state.showTemplate,
    });
    if (this.state.showTemplate) {
      const { dispatch } = this.props;
      dispatch({ 
        type: 'deploy/getDeploy',
        payload: {
          page: 1,
          pageSize: 10,
        }
      });
    }
  };

  switchUndoTemplate = () => {
    this.setState({ 
      undoTemplate: !this.state.undoTemplate,
    });
    if (this.state.undoTemplate) {
      const { dispatch } = this.props;
      dispatch({ 
        type: 'deploy/getDeploy',
        payload: {
          page: 1,
          pageSize: 10,
        }
      });
    }
  };

  handleEnvChange = (values) => {
    if (values) {
      const { dispatch } = this.props;
      dispatch({
        type: 'config/getProject',
        payload: {
          page: 1,
          pageSize: 999,
          active: 1,
          envId: values,
        }
      });
    }
  };

  handleAppChange = (values) => {
    if (values) {
      const { dispatch } = this.props;
      dispatch({
        type: 'config/getAppTemplate',
        payload: {
          aid: values.id,
        }
      });
    }

    this.setState({ 
      selectAppVisible: false,
      selectTemplateVisible:true,
    });
  };

  handleExt1Ok = () => {
    this.setState({ 
      visible: true,
      selectAppVisible: false,
    });
  };

  showExt1Form = (projectId) => {
    // 获取该项目的所有模板

    this.setState({ 
      selectAppVisible: false,
      ext1Visible:true,
    });
  }

  showExt2Form = (projectId) => {
    // 获取该项目的所有模板

    this.setState({ 
      selectAppVisible: false,
      ext2Visible:true,
    });
  }

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'deploy/deployDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  }

  // Popconfirm 取消事件
  cancel = () => {
  }

  // Check Undo time
  undoCheck = (value) => {
    if (moment().diff(moment(value)) > 86400000) {
      return true
    }
    return false
  }

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑-' + values.Name;
    const { dispatch } = this.props;
    dispatch({
      type: 'config/getAppTemplate',
      payload: {
        aid: values.Aid,
      }
    }).then(()=>{
      this.setState({ 
        editCacheData: values,
        selectTemplateVisible: true,
      })
    })
  };

  handleDeployYes = (values, isRollBack=false, isLog=false) => {
    const list = this.props.appTemplateList
    var tmpExtend = 0
    var type = 1
    var log  = 0
    for (var i = list.length - 1; i >= 0; i--) {
      if (list[i]['Dtid'] === values.Tid) {
        tmpExtend = list[i]['Extend']
        break
      }
    }
    if (isRollBack) {
      type = 2
    }
    if (isLog) {
      log = 1
    }
    if (tmpExtend == 1) {
      return "/deploy/do/common?type=" + type + "&id=" + values.ID + "&log=" + log + "&templateName=" + values.TemplateName
    } else {
      return "/deploy/do/custom?type=" + type + "&id=" + values.ID + "&log=" + log + "&templateName=" + values.TemplateName
    }
  };

  handleRollback = (info) => {
    this.setState({loading: true});
    httpGet(`/admin/undo/request/${info.ID}`).then(res => {
      if (res.code != 200) {
        message.error(res.message)
        return
      }

      if (res.code == 200) {
         Modal.confirm({
          title: '回滚确认',
          content: `确定要回滚至 ${timeDatetimeTrans(res.data['UpdateTime'])} 发布的名称为【${res.data['Version']}】的发布申请版本?`,
          onOk: () => {this.handleRollbackConfirm(info.ID, res.data['Version'])}
        })
      }
    }).finally(() => this.setState({loading: false}))
  }

  handleRollbackConfirm = (id, version) => {
    var Params = {id:id, Version: version}
    const { dispatch } = this.props;
    dispatch({ 
      type: 'deploy/rollbackConfirm',
      payload: Params
    })
  }

  undeployYes = (values) => {
    this.setState({ 
      undoTemplate: true ,
      id: values,
    });
  };

  // 发布
  deployYes = (values) => {
    this.setState({ 
      showTemplate: true ,
      id: values,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/getDeploy',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  columns = [
    {
      title: '申请标题',
      dataIndex: 'Name',
    },
    {
      title: '上线模板',
      dataIndex: 'Tid',
      'render': Tid => this.props.appTemplateList.map(x => {
        if (Tid==x.Dtid) {
          return x.TemplateName
        }
      })
    },
    {
      title: '分支/Tag',
      dataIndex: 'TagBranch',
    },
    {
      title: '状态',
      dataIndex: 'Status',
      ellipsis: true,
      'render': (Status, info) => {
        if (Status == 1) {
          return <Tag>待审核</Tag>
        } else if (Status == 3) {
          return <Tag color="blue">待发布</Tag>
        } else if (Status == 4) {
          return <Tag color="blue">回滚待发布</Tag>
        } else if (Status == 5) {
          return <Tag color="green">发布成功</Tag>
        } else if (Status == 6) {
          return <Tag color="green">回滚成功</Tag>
        } else if (Status == -1) {
          return  <Popover title="驳回意见:" content={info.Reason}>
                    <span style={{color: '#1890ff'}}>已驳回</span>
                  </Popover>
        } else if (Status == -2) {
          return <Tag color="red">回滚失败</Tag>
        } else if (Status == -3) {
          return <Tag color="red">发布异常</Tag>
        } else if (Status == 2) {
          if (info.Reason != "") {
            return <Popover title="审核意见:" content={info.Reason}>
                      <span style={{color: '#1890ff'}}>待发布</span>
                    </Popover>
          } else {
            return <Tag color="blue">待发布</Tag>
          }
        }
      }
    },
    {
      title: '状态变更时间',
      dataIndex: 'UpdateTime',
      'render': UpdateTime => timeDatetimeTrans(UpdateTime),
    },
    {
      title: '操作',
      key: 'action',
      render: (text, record) => (
        <span>
          {  
            hasPermission('deploy-app-request') && record.Status == 2   &&
            // <a title="发布" onClick={()=>this.handleDeployYes(record.ID)}><Icon type="redo" />发布</a>
            <Link to={this.handleDeployYes(record)}>
              <Icon type="redo"/>发布
            </Link>
          }
          {  
            hasPermission('deploy-app-request') &&  record.Status == 4 &&
            // <a title="发布" onClick={()=>this.handleDeployYes(record.ID)}><Icon type="redo" />发布</a>
            <Link to={this.handleDeployYes(record, true, false)}>
              <Icon type="redo"/>发布
            </Link>
          }
          {hasPermission('deploy-app-undo') && record.Status == 5  &&
            <a disabled={this.undoCheck(record.UpdateTime)} onClick={()=>{this.handleRollback(record)}} title="回滚" >
              <Icon type="undo" />回滚
            </a>
          }
          {hasPermission('deploy-app-review') && record.Status == 1  &&
            <a title="删除" onClick={()=>this.handleApprove(record.ID)}><Icon type="team" />审核</a>
          }
          { record.Status == 1  &&  <Divider type="vertical" /> }
          {hasPermission('deploy-app-edit') && record.Status == 1  && 
            <a onClick={()=>{this.handleEdit(record)}}>
              <Icon type="edit"/>编辑
            </a>
          }
          { record.Status <= 2 && <Divider type="vertical" /> }
          {hasPermission('deploy-app-del') && record.Status <= 2  &&
            <Popconfirm title="你确定要删除吗?"
              onConfirm={()=>{this.deleteRecord(record.ID)}} onCancel={()=>{this.cancel()}}>
              <a title="删除" ><Icon type="delete" />删除</a>
            </Popconfirm>
          }

          { (record.Status >= 5 || record.Status < -1) && <Divider type="vertical" /> }
          {hasPermission('deploy-app-view') && record.Status == 5  &&
            <Link to={this.handleDeployYes(record, false, true)}>
              <Icon type="eye"/>查看
            </Link>
          }
          {hasPermission('deploy-app-view') && record.Status == 6  &&
            <Link to={this.handleDeployYes(record, false, true)}>
              <Icon type="eye"/>查看
            </Link>
          }
          {hasPermission('deploy-app-view') && record.Status < -1  &&
            <Link to={this.handleDeployYes(record, false, true)}>
              <Icon type="eye"/>查看
            </Link>
          }
        </span>
      ),
    },
  ];

  render() {
    const {selectAppVisible, showTemplate, undoTemplate, gitShow,
      approveVisible, selectTemplateVisible,
      disableValues, commitLoading, branchLoading, ext1Visible, ext2Visible} = this.state;

    const {deployList, deployListLoading, deployLen, projectsList,
      configEnvList, gitBranchList, gitCommitList, output,
      appTemplateList, form: { getFieldDecorator } } = this.props;

    const addvar = <Button type="" onClick={this.showAddModal} >提单</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('deploy-app-newline') && <Col span={10}>{addvar}</Col>}
      </Row>;
    return (
      <div>
        {selectAppVisible && 
          <SelectApp 
            configEnvList={configEnvList}
            projectsList={projectsList}
            onCancel={this.selectCancel}
            showExt1Form={this.showExt1Form}
            showExt2Form={this.showExt2Form}
            handleAppChange={this.handleAppChange}
          />
        }
        {selectTemplateVisible && 
          <SelectTemplate 
            onCancel={this.selectCancel}
            appTemplateList={appTemplateList}
            editCacheData={this.state.editCacheData}
          />
        }
        {ext1Visible && <Ext1Form/>}
        {ext2Visible && <Ext2Form/>}
        {approveVisible && 
          <Approve 
            approveCanael={this.approveCanael}
            id={this.state.id}
            dispatch={this.props.dispatch}
          />
        }

        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: deployLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={deployList} loading={deployListLoading} rowKey="ID" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(DeployListPage);
